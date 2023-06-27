package winsigncert

/*
#cgo LDFLAGS: -lcrypt32 -lwintrust
#cgo CFLAGS: -w
#include <windows.h>
#include <stdio.h>
#include <wincrypt.h>
#include <wintrust.h>
#include <Softpub.h>

#define nullptr NULL
#define false FALSE
#define true TRUE

#define NORMAL_SIZE 200
#define SUBJECT_SIZE 400
#define TIMESTAMP_SIZE 50
#define LONG_MAX_PATH 2048

#define ENCODING (X509_ASN_ENCODING | PKCS_7_ASN_ENCODING)

typedef struct {
	LPWSTR lpszProgramName;
	LPWSTR lpszPublisherLink;
	LPWSTR lpszMoreInfoLink;
} SPROG_PUBLISHERINFO, *PSPROG_PUBLISHERINFO;

LPWSTR AllocateAndCopyWideString(LPCWSTR inputString)
{
	LPWSTR outputString = NULL;

	outputString = (LPWSTR)LocalAlloc(LPTR, (wcslen(inputString) + 1) * sizeof(WCHAR));
	if (outputString != NULL)
		lstrcpyW(outputString, inputString);

	return outputString;
}

BOOL GetTimeStampSignerInfo(PCMSG_SIGNER_INFO pSignerInfo, PCMSG_SIGNER_INFO *pCounterSignerInfo)
{
	PCCERT_CONTEXT pCertContext = NULL;
	BOOL fReturn = FALSE;
	BOOL fResult;
	DWORD dwSize;

	if (pSignerInfo != NULL | pCounterSignerInfo != NULL)
	{
		*pCounterSignerInfo = NULL;

		// Loop through unathenticated attributes for
		// szOID_RSA_counterSign OID.
		for (DWORD n = 0; n < pSignerInfo->UnauthAttrs.cAttr; n++)
		{
			if (lstrcmpA(pSignerInfo->UnauthAttrs.rgAttr[n].pszObjId,
				szOID_RSA_counterSign) == 0)
			{
				// Get size of CMSG_SIGNER_INFO structure.
				fResult = CryptDecodeObject(ENCODING,
					PKCS7_SIGNER_INFO,
					pSignerInfo->UnauthAttrs.rgAttr[n].rgValue[0].pbData,
					pSignerInfo->UnauthAttrs.rgAttr[n].rgValue[0].cbData,
					0,
					NULL,
					&dwSize);
				if (!fResult)
					return false;

				// Allocate memory for CMSG_SIGNER_INFO.
				*pCounterSignerInfo = (PCMSG_SIGNER_INFO)LocalAlloc(LPTR, dwSize);
				if (!*pCounterSignerInfo)
					return false;
			}

			// Decode and get CMSG_SIGNER_INFO structure
			// for timestamp certificate.
			fResult = CryptDecodeObject(ENCODING,
				PKCS7_SIGNER_INFO,
				pSignerInfo->UnauthAttrs.rgAttr[n].rgValue[0].pbData,
				pSignerInfo->UnauthAttrs.rgAttr[n].rgValue[0].cbData,
				0,
				(PVOID)*pCounterSignerInfo,
				&dwSize);
			if (!fResult)
				return false;

			fReturn = TRUE;

			break; // Break from for loop.
		}
	}
	return fReturn;
}

BOOL GetDateOfTimeStamp(PCMSG_SIGNER_INFO pSignerInfo, SYSTEMTIME *st)
{
	BOOL fResult;
	FILETIME lft, ft;
	DWORD dwData;
	BOOL fReturn = FALSE;

	// Loop through authenticated attributes and find
	// szOID_RSA_signingTime OID.
	for (DWORD n = 0; n < pSignerInfo->AuthAttrs.cAttr; n++)
	{
		if (lstrcmpA(szOID_RSA_signingTime,
			pSignerInfo->AuthAttrs.rgAttr[n].pszObjId) == 0)
		{
			// Decode and get FILETIME structure.
			dwData = sizeof(ft);
			fResult = CryptDecodeObject(ENCODING,
				szOID_RSA_signingTime,
				pSignerInfo->AuthAttrs.rgAttr[n].rgValue[0].pbData,
				pSignerInfo->AuthAttrs.rgAttr[n].rgValue[0].cbData,
				0,
				(PVOID)&ft,
				&dwData);
			if (!fResult)
			{
				//printf(("CryptDecodeObject failed with %x\n"),
				//	GetLastError());
				break;
			}

			// Convert to local time.
			FileTimeToLocalFileTime(&ft, &lft);
			FileTimeToSystemTime(&lft, st);

			fReturn = TRUE;

			break; // Break from for loop.

		} //lstrcmp szOID_RSA_signingTime
	} // for

	return fReturn;
}

BOOL GetCertificateInfo(PCCERT_CONTEXT pCertContext, LPTSTR* pszName, LPTSTR* psubject)
{
	BOOL bResult = false;
	LPTSTR szName = NULL;
	LPTSTR subject = NULL;
	DWORD dwData;

	if (!pCertContext) {
		return false;
	}

	dwData = pCertContext->pCertInfo->SerialNumber.cbData;
	if (!dwData) {
		return false;
	}

	//for (DWORD n = 0; n < dwData; n++)
	//	pCertContext->pCertInfo->SerialNumber.pbData[dwData - (n + 1)];

	// Get Issuer name.
	if (!(dwData = CertGetNameString(pCertContext, CERT_NAME_SIMPLE_DISPLAY_TYPE, CERT_NAME_ISSUER_FLAG, NULL, NULL, 0)))
		goto Exit0;

	szName = (LPTSTR)LocalAlloc(LPTR, dwData * sizeof(TCHAR));
	if (!szName)
		goto Exit0;

	if (!(CertGetNameString(pCertContext, CERT_NAME_SIMPLE_DISPLAY_TYPE, CERT_NAME_ISSUER_FLAG, NULL, szName, dwData)))
		goto Exit0;

	// Get Subject name size.
	if (!(dwData = CertGetNameString(pCertContext, CERT_NAME_SIMPLE_DISPLAY_TYPE, 0, NULL, NULL, 0)))
		goto Exit0;

	subject = (LPTSTR)LocalAlloc(LPTR, dwData * sizeof(TCHAR));
	if (!subject)
		goto Exit0;

	// Get subject name.
	if (!(CertGetNameString(pCertContext, CERT_NAME_SIMPLE_DISPLAY_TYPE, 0, NULL, subject, dwData)))
		goto Exit0;

	*pszName = szName;
	*psubject = subject;
	bResult = true;
Exit0:
	if (!bResult){
		if (szName)
			LocalFree(szName);

		if (subject)
			LocalFree(subject);
	}

	return bResult;
}

BOOL GetProgAndPublisherInfo(PCMSG_SIGNER_INFO pSignerInfo, PSPROG_PUBLISHERINFO Info)
{
	BOOL bResult = FALSE;
	PSPC_SP_OPUS_INFO OpusInfo = NULL;
	DWORD dwData;
	BOOL fResult;

	if (!pSignerInfo || !Info)
	{
		return false;
	}

	// Loop through authenticated attributes and find
	// SPC_SP_OPUS_INFO_OBJID OID.
	for (DWORD n = 0; n < pSignerInfo->AuthAttrs.cAttr; n++)
	{
		if (lstrcmpA(SPC_SP_OPUS_INFO_OBJID,
			pSignerInfo->AuthAttrs.rgAttr[n].pszObjId) != 0)
		{
			continue;
		}

		// Get Size of SPC_SP_OPUS_INFO structure.
		fResult = CryptDecodeObject(ENCODING,
			SPC_SP_OPUS_INFO_OBJID,
			pSignerInfo->AuthAttrs.rgAttr[n].rgValue[0].pbData,
			pSignerInfo->AuthAttrs.rgAttr[n].rgValue[0].cbData,
			0,
			NULL,
			&dwData);
		if (!fResult)
			goto Exit0;

		// Allocate memory for SPC_SP_OPUS_INFO structure.
		OpusInfo = (PSPC_SP_OPUS_INFO)LocalAlloc(LPTR, dwData);
		if (!OpusInfo)
			goto Exit0;

		// Decode and get SPC_SP_OPUS_INFO structure.
		fResult = CryptDecodeObject(ENCODING,
			SPC_SP_OPUS_INFO_OBJID,
			pSignerInfo->AuthAttrs.rgAttr[n].rgValue[0].pbData,
			pSignerInfo->AuthAttrs.rgAttr[n].rgValue[0].cbData,
			0,
			OpusInfo,
			&dwData);
		if (!fResult)
			goto Exit0;

		// Fill in Program Name if present.
		if (OpusInfo->pwszProgramName)
		{
			Info->lpszProgramName = (LPWSTR)AllocateAndCopyWideString(OpusInfo->pwszProgramName);
		}
		else
			Info->lpszProgramName = NULL;

		// Fill in Publisher Information if present.
		if (OpusInfo->pPublisherInfo)
		{

			switch (OpusInfo->pPublisherInfo->dwLinkChoice)
			{
			case SPC_URL_LINK_CHOICE:
				Info->lpszPublisherLink =
					AllocateAndCopyWideString(OpusInfo->pPublisherInfo->pwszUrl);
				break;

			case SPC_FILE_LINK_CHOICE:
				Info->lpszPublisherLink =
					AllocateAndCopyWideString(OpusInfo->pPublisherInfo->pwszFile);
				break;

			default:
				Info->lpszPublisherLink = NULL;
				break;
			}
		}
		else
		{
			Info->lpszPublisherLink = NULL;
		}

		// Fill in More Info if present.
		if (OpusInfo->pMoreInfo)
		{
			switch (OpusInfo->pMoreInfo->dwLinkChoice)
			{
			case SPC_URL_LINK_CHOICE:
				Info->lpszMoreInfoLink =
					AllocateAndCopyWideString(OpusInfo->pMoreInfo->pwszUrl);
				break;

			case SPC_FILE_LINK_CHOICE:
				Info->lpszMoreInfoLink =
					AllocateAndCopyWideString(OpusInfo->pMoreInfo->pwszFile);
				break;

			default:
				Info->lpszMoreInfoLink = NULL;
				break;
			}
		}
		else
		{
			Info->lpszMoreInfoLink = NULL;
		}
		goto Exit1;
	} // for

Exit1:
	bResult = true;
Exit0:
	return bResult;
}

typedef struct DEPTINFO {
	char* publisher;
	char* programName;
	char* subject;
	char* timestamp;
} DEPTINFO;

DEPTINFO* GatherInfo(char* path)
{
	WCHAR szFileName[LONG_MAX_PATH];
	HCERTSTORE hStore = NULL;
	HCRYPTMSG hMsg = NULL;
	PCCERT_CONTEXT pCertContext = NULL;
	BOOL fResult;
	DWORD dwEncoding, dwContentType, dwFormatType;
	PCMSG_SIGNER_INFO pSignerInfo = NULL;
	PCMSG_SIGNER_INFO pCounterSignerInfo = NULL;
	DWORD dwSignerInfo;
	CERT_INFO CertInfo;
	SPROG_PUBLISHERINFO ProgPubInfo;
	SYSTEMTIME st;

	DEPTINFO* info = (DEPTINFO*)malloc(sizeof(DEPTINFO));

	memset(info, 0x00, sizeof(DEPTINFO));
	info->programName = (char*)malloc(NORMAL_SIZE);
	info->publisher = (char*)malloc(NORMAL_SIZE);
	info->subject = (char*)malloc(SUBJECT_SIZE);
	info->timestamp = (char*)malloc(TIMESTAMP_SIZE);

	ZeroMemory(info->programName, NORMAL_SIZE);
	ZeroMemory(info->publisher, NORMAL_SIZE);
	ZeroMemory(info->subject, SUBJECT_SIZE);
	ZeroMemory(info->timestamp, TIMESTAMP_SIZE);
	ZeroMemory(&ProgPubInfo, sizeof(ProgPubInfo));

	if (!path || strlen(path) <= 1) {
		return info;
	}

	if (mbstowcs(szFileName, path, LONG_MAX_PATH) == -1)
		return info;

	fResult = CryptQueryObject(CERT_QUERY_OBJECT_FILE, szFileName, CERT_QUERY_CONTENT_FLAG_PKCS7_SIGNED_EMBED,
		CERT_QUERY_FORMAT_FLAG_BINARY, 0,
		&dwEncoding,
		&dwContentType,
		&dwFormatType,
		&hStore,
		&hMsg,
		NULL);
	if (!fResult)
		goto Exit0;

	fResult = CryptMsgGetParam(hMsg, CMSG_SIGNER_INFO_PARAM, 0, NULL, &dwSignerInfo);
	if (!fResult)
		goto Exit0;

	pSignerInfo = (PCMSG_SIGNER_INFO)LocalAlloc(LPTR, dwSignerInfo);
	if (!pSignerInfo)
		goto Exit0;

	// Get Signer Information.
	fResult = CryptMsgGetParam(hMsg, CMSG_SIGNER_INFO_PARAM, 0, (PVOID)pSignerInfo, &dwSignerInfo);
	if (!fResult)
		goto Exit0;

	if (GetProgAndPublisherInfo(pSignerInfo, &ProgPubInfo))
	{
		if (ProgPubInfo.lpszProgramName != NULL)
		{
			wcstombs(info->programName, ProgPubInfo.lpszProgramName, NORMAL_SIZE);
		}
	}

	CertInfo.Issuer = pSignerInfo->Issuer;
	CertInfo.SerialNumber = pSignerInfo->SerialNumber;

	pCertContext = CertFindCertificateInStore(hStore,
		ENCODING,
		0,
		CERT_FIND_SUBJECT_CERT,
		(PVOID)&CertInfo,
		NULL);
	if (!pCertContext)
		goto Exit0;

	if(!GetCertificateInfo(pCertContext, &info->publisher, &info->subject))
		goto Exit0;

	if (GetTimeStampSignerInfo(pSignerInfo, &pCounterSignerInfo))
	{
		// Search for Timestamp certificate in the temporary
		// certificate store.
		CertInfo.Issuer = pCounterSignerInfo->Issuer;
		CertInfo.SerialNumber = pCounterSignerInfo->SerialNumber;

		pCertContext = CertFindCertificateInStore(hStore,
			ENCODING,
			0,
			CERT_FIND_SUBJECT_CERT,
			(PVOID)&CertInfo,
			NULL);
		if (!pCertContext)
			goto Exit0;

		//PrintCertificateInfo(pCertContext);

		if (GetDateOfTimeStamp(pCounterSignerInfo, &st))
		{
			snprintf(info->timestamp, TIMESTAMP_SIZE, "%04d-%02d-%02d %02d:%02d:%02d",
				st.wYear,
				st.wMonth,
				st.wDay,
				st.wHour,
				st.wMinute,
				st.wSecond);
		}
	}

Exit0:
	// CLEANING
	if (ProgPubInfo.lpszProgramName != NULL)
		LocalFree(ProgPubInfo.lpszProgramName);
	if (ProgPubInfo.lpszPublisherLink != NULL)
		LocalFree(ProgPubInfo.lpszPublisherLink);
	if (ProgPubInfo.lpszMoreInfoLink != NULL)
		LocalFree(ProgPubInfo.lpszMoreInfoLink);

	if (pSignerInfo != NULL) LocalFree(pSignerInfo);
	if (pCounterSignerInfo != NULL) LocalFree(pCounterSignerInfo);
	if (pCertContext != NULL) CertFreeCertificateContext(pCertContext);
	if (hStore != NULL) CertCloseStore(hStore, 0);
	if (hMsg != NULL) CryptMsgClose(hMsg);

	return info;
}

int VerifyEmbeddedSignature(LPCWSTR pwszSourceFile)
{
	int nResult = 0;
	LONG lStatus;
	DWORD dwLastError;

	// Initialize the WINTRUST_FILE_INFO structure.

	WINTRUST_FILE_INFO FileData;
	memset(&FileData, 0, sizeof(FileData));
	FileData.cbStruct = sizeof(WINTRUST_FILE_INFO);
	FileData.pcwszFilePath = pwszSourceFile;
	FileData.hFile = NULL;
	FileData.pgKnownSubject = NULL;

	GUID WVTPolicyGUID = WINTRUST_ACTION_GENERIC_VERIFY_V2;
	WINTRUST_DATA WinTrustData;

	// Initialize the WinVerifyTrust input data structure.

	// Default all fields to 0.
	memset(&WinTrustData, 0, sizeof(WinTrustData));

	WinTrustData.cbStruct = sizeof(WinTrustData);

	// Use default code signing EKU.
	WinTrustData.pPolicyCallbackData = NULL;

	// No data to pass to SIP.
	WinTrustData.pSIPClientData = NULL;

	// Disable WVT UI.
	WinTrustData.dwUIChoice = WTD_UI_NONE;

	// No revocation checking.
	WinTrustData.fdwRevocationChecks = WTD_REVOKE_NONE;

	// Verify an embedded signature on a file.
	WinTrustData.dwUnionChoice = WTD_CHOICE_FILE;

	// Verify action.
	WinTrustData.dwStateAction = WTD_STATEACTION_VERIFY;

	// Verification sets this value.
	WinTrustData.hWVTStateData = NULL;

	// Not used.
	WinTrustData.pwszURLReference = NULL;

	// This is not applicable if there is no UI because it changes
	// the UI to accommodate running applications instead of
	// installing applications.
	WinTrustData.dwUIContext = 0;

	// Set pFile.
	WinTrustData.pFile = &FileData;

	// WinVerifyTrust verifies signatures as specified by the GUID
	// and Wintrust_Data.
	lStatus = WinVerifyTrust(
		NULL,
		&WVTPolicyGUID,
		&WinTrustData);

	switch (lStatus)
	{
	case ERROR_SUCCESS:
		break;
	case TRUST_E_NOSIGNATURE:
	case TRUST_E_EXPLICIT_DISTRUST:
	case TRUST_E_SUBJECT_NOT_TRUSTED:
	case CRYPT_E_SECURITY_SETTINGS:
	default:
		goto Exit0;
	}

	nResult = 1;
Exit0:
	// Any hWVTStateData must be released by a call with close.
	WinTrustData.dwStateAction = WTD_STATEACTION_CLOSE;

	lStatus = WinVerifyTrust(
		NULL,
		&WVTPolicyGUID,
		&WinTrustData);

	return nResult;
}
*/
import "C"
import (
	"unicode/utf16"
	"unsafe"
)

type DEPTINFO struct {
	ProgramName *string
	Publisher   *string
	Subject     *string
	Timestamp   *string
}

func encode(s string) C.LPCWSTR {
	wstr := utf16.Encode([]rune(s))
	wstr = append(wstr, 0x00)
	return (C.LPCWSTR)(unsafe.Pointer(&wstr[0]))
}

func ValidateSignCert(winFilePath string) bool {
	pathCString := encode(winFilePath)
	resp := C.VerifyEmbeddedSignature(pathCString)
	if resp == 1 {
		return true
	}
	return false
}

func GetSignCertInfo(winFilePath string) *DEPTINFO {
	pathCString := C.CString(winFilePath)
	info := C.GatherInfo(pathCString)

	name := C.GoString(info.programName)
	publisher := C.GoString(info.publisher)
	suby := C.GoString(info.subject)
	timestamp := C.GoString(info.timestamp)
	C.free(unsafe.Pointer(info))

	return &DEPTINFO{
		ProgramName: &name,
		Publisher:   &publisher,
		Subject:     &suby,
		Timestamp:   &timestamp,
	}
}
