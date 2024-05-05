logLevel: debug
logs:
  - filename: "goutils.log"
    stdout: true
    fileOut: true
    httpOut: false
    outputEncoder: console 
    httpUrl: "http://127.0.0.1:8053/goutils/log"    
    httpDebug: false    
  - filename: "goutils.json"
    stdout: true
    fileOut: true
    #httpOut: true
    outputEncoder: json
    #httpUrl: "http://127.0.0.1:8053/goutils/log"
    #httpDebug: true
