package checksum

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils"
)

const OS_FILE_R_W = 0666

type ChecksumInfo struct {
	FilePath string
	Crc32Val string
	FileSize string
}

func ChecksumCrc32(fileReader io.Reader) (uint32, error) {
	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, fileReader); err != nil {
		return 0, err
	}
	return hash.Sum32(), nil
}

// AddFolderSuffix 为路径添加分隔符后缀
func AddFolderSuffix(folder string) string {
	separator := string(os.PathSeparator)
	if strings.HasSuffix(folder, separator) {
		return folder
	}
	return folder + separator
}

// CompareChecksumFiles
// srcDir 校验目标文件夹
// checksumPath 校验文件路径
// ignores 忽略的文件
func CompareChecksumFiles(ctx context.Context, root string, checksumPath string) error {
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return errors.New(fmt.Sprintf("open src dir %s faild", root))
	}
	checkInfo, err := GetChecksumInfo(ctx, checksumPath)
	if err != nil {
		return err
	}
	paths := GetCheckFileList(ctx, checkInfo)
	if err != nil {
		return err
	}
	err = ChecksumFilesWithCheckInfo(root, checkInfo, paths)
	if err != nil {
		return err
	}
	return nil
}

func IsChecksumFileValid(ctx context.Context, checksumPath, md5Path string) bool {
	bContent, err := ioutil.ReadFile(md5Path)
	if err != nil {
		log.Error(ctx, "open md5Path:%s error:%v", md5Path, err)
		return false
	}
	fd, err := os.Open(checksumPath)
	if err != nil {
		log.Error(ctx, "open checksumPath:%s error:%v", checksumPath, err)
		return false
	}
	defer fd.Close()
	md5Handle := md5.New()
	_, err = io.Copy(md5Handle, fd)
	if err != nil {
		log.Error(ctx, "checksumPath:%s md5 error:%v", checksumPath, err)
		return false
	}
	md5Val := md5Handle.Sum(nil)
	hexMd5 := hex.EncodeToString(md5Val)
	if bytes.Equal(bContent, []byte(hexMd5)) {
		return true
	}
	log.Error(ctx, "checksumPath:%s md5 is not equal with md5 file:%s", checksumPath, md5Path)
	return false
}

func ChecksumFilesWithCheckInfo(root string, checkInfo map[string]*ChecksumInfo, files []string) error {
	for _, file := range files {
		fileInfo, err := GetFileInfo(root, file)
		if err != nil {
			return err
		}
		orginInfo, ok := checkInfo[file]
		if !ok || orginInfo == nil {
			return errors.New(fmt.Sprintf("filename:%s orgin check info invalid", file))
		}
		if (fileInfo.Crc32Val != orginInfo.Crc32Val) || (fileInfo.FileSize != orginInfo.FileSize) {
			return errors.New(fmt.Sprintf("filename:%s checksum not equal", file))
		}
	}
	return nil
}

func GetFileInfo(root, fileName string) (*ChecksumInfo, error) {
	filePath := filepath.Join(root, fileName)
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(filePath) //path是相对路径
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	value, _ := ChecksumCrc32(fd)
	return &ChecksumInfo{
		FilePath: fileName,
		Crc32Val: strconv.FormatUint(uint64(value), 10),
		FileSize: strconv.FormatUint(uint64(info.Size()), 10),
	}, nil
}

// GenerateChecksumFile 生成checksum文件
func GenerateChecksumFile(ctx context.Context, folder string, checksumName string) (checkSumPath string, err error) {
	return GenerateChecksumFileWithIgnore(ctx, folder, checksumName, []string{})
}

// GenerateChecksumFileWithIgnore 在排除某些文件的基础上，生成checksum
func GenerateChecksumFileWithIgnore(ctx context.Context, folder string, checksumName string, ignores []string) (checkSumPath string, err error) {
	if folder == "" || checksumName == "" {
		err = errors.New("folder or checksum can`t be empty")
		return
	}

	folder = AddFolderSuffix(folder)
	paths, err := PopulateFilePathsRecursively(ctx, folder, ignores)
	if err != nil {
		log.Error(ctx, "open folder failed. err:%v", err)
		return
	}
	builder := strings.Builder{}
	for idx, path := range paths {
		file, err := os.Open(filepath.Join(folder, path))
		if err != nil {
			log.Error(ctx, "open filepath:%s failed. err:%v", path, err)
			return "", err
		}
		checksumCrc32, err := ChecksumCrc32(file)
		if err != nil {
			log.Error(ctx, "get filepath:%s crc32 val failed. err:%v", path, err)
			return "", err
		}
		fileInfo, _ := file.Stat()
		if idx == len(paths)-1 {
			builder.WriteString(fmt.Sprintf("%v|%v|%v", strings.ReplaceAll(path, string(os.PathSeparator), "/"), checksumCrc32, fileInfo.Size()))
		} else {
			builder.WriteString(fmt.Sprintf("%v|%v|%v\n", strings.ReplaceAll(path, string(os.PathSeparator), "/"), checksumCrc32, fileInfo.Size()))
		}
	}
	// 生成文件checksum文件
	checkSumPath = filepath.Join(folder, fmt.Sprintf("%s.checksum", checksumName))
	err = ioutil.WriteFile(checkSumPath, []byte(builder.String()), OS_FILE_R_W)
	return
}

// GenerateChecksumMd5File 生成checksum.md5文件
func GenerateChecksumMd5File(ctx context.Context, checksumPath string) (checksumMd5Path string, err error) {
	//打开文件
	file, err := os.Open(checksumPath)
	if err != nil {
		log.Error(ctx, "open checksumPath:%s error:%v", checksumPath, err)
		return
	}
	defer file.Close()
	// 生成md5
	md5Handle := md5.New()
	_, err = io.Copy(md5Handle, file)
	if err != nil {
		log.Error(ctx, "copy file failed, err：%v", err)
		return
	}
	md5Val := md5Handle.Sum(nil)
	hexVal := hex.EncodeToString(md5Val)

	log.Info(ctx, "checksumPath:%s generate md5 val: %s", checksumPath, hexVal)
	checksumMd5Path = checksumPath + ".md5"
	// 写入文件
	err = ioutil.WriteFile(checksumMd5Path, []byte(hexVal), OS_FILE_R_W)
	return
}

// PopulateFilePathsRecursively 递归获取文件夹内所有文件的路径
func PopulateFilePathsRecursively(ctx context.Context, folder string, ignores []string) ([]string, error) {
	paths := make([]string, 0)
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error(ctx, "folder [%v] scan meet some err: %v, should skip", folder, err)
			return err
		}
		relPath, _ := filepath.Rel(folder, path)
		isContain, _ := utils.StringsInArray(ignores, relPath)
		if info.IsDir() && isContain {
			return filepath.SkipDir
		}
		if info.IsDir() || isContain {
			return nil
		}

		paths = append(paths, relPath)
		log.Debug(ctx, "folder: %v", info.Name())
		return nil
	})
	return paths, nil
}

func GetChecksumInfo(ctx context.Context, checksumPath string) (checkInfo map[string]*ChecksumInfo, err error) {
	checkInfo = make(map[string]*ChecksumInfo)
	_, err = os.Stat(checksumPath)
	if err != nil {
		return
	}
	content, err := ioutil.ReadFile(checksumPath)
	if err != nil {
		log.Error(ctx, "read checksum file error:%v", err)
		return checkInfo, err
	}
	reader := bufio.NewReader(bytes.NewReader(content))
	for {
		line, _, lineErr := reader.ReadLine()
		if lineErr == io.EOF {
			break
		}
		spVals := strings.Split(string(line), "|")

		if len(spVals) == 3 {
			checkInfo[spVals[0]] = &ChecksumInfo{
				FilePath: spVals[0],
				Crc32Val: spVals[1],
				FileSize: spVals[2],
			}
		} else {
			return checkInfo, errors.New("checksum file content invalid")
		}
	}
	return checkInfo, nil
}

func GetCheckFileList(ctx context.Context, checkInfo map[string]*ChecksumInfo) []string {
	var files []string
	for k, _ := range checkInfo {
		files = append(files, k)
	}
	return files
}

func RelWalkInfo(ctx context.Context, root string, ignores ...string) ([]string, error) {
	paths, err := WalkInfo(ctx, root, ignores...)
	if err != nil {
		return nil, err
	}
	var relPaths []string
	for _, path := range paths {
		relPath, _ := filepath.Rel(root, path)
		relPaths = append(relPaths, strings.ReplaceAll(relPath, string(os.PathSeparator), "/"))
	}
	return relPaths, nil
}

func WalkInfo(ctx context.Context, root string, ignores ...string) ([]string, error) {
	var err error
	var fds []os.FileInfo
	var files []string
	if _, err = os.Lstat(root); err != nil {
		return nil, err
	}
	if fds, err = ioutil.ReadDir(root); err != nil {
		return nil, err
	}

	for _, fd := range fds {
		if ok, _ := utils.StringsInArray(ignores, fd.Name()); ok {
			log.Debug(context.Background(), "walk info ignore file path:%v", fd.Name())
			continue
		}
		if fd.IsDir() {
			subFiles, err := WalkInfo(ctx, root, ignores...)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, filepath.Join(root, fd.Name()))
		}
	}
	return files, nil
}
