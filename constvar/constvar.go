package constvar

import (
	"fmt"
	"time"
)

const (
	// API_VERSION = "0.3.2"
	// http://www.network-science.de/ascii/
	// http://patorjk.com/software/taag
	// https://www.jianshu.com/p/fca56d635091
	LOGO_ASCII = `
                                            
                                            
     QQQQQQQQQ      NNNNNNNN        NNNNNNNN
   QQ:::::::::QQ    N:::::::N       N::::::N
 QQ:::::::::::::QQ  N::::::::N      N::::::N
Q:::::::QQQ:::::::Q N:::::::::N     N::::::N
Q::::::O   Q::::::Q N::::::::::N    N::::::N
Q:::::O     Q:::::Q N:::::::::::N   N::::::N
Q:::::O     Q:::::Q N:::::::N::::N  N::::::N
Q:::::O     Q:::::Q N::::::N N::::N N::::::N
Q:::::O     Q:::::Q N::::::N  N::::N:::::::N
Q:::::O     Q:::::Q N::::::N   N:::::::::::N
Q:::::O  QQQQ:::::Q N::::::N    N::::::::::N
Q::::::O Q::::::::Q N::::::N     N:::::::::N
Q:::::::QQ::::::::Q N::::::N      N::::::::N
 QQ::::::::::::::Q  N::::::N       N:::::::N
   QQ:::::::::::Q   N::::::N        N::::::N
     QQQQQQQQ::::QQ NNNNNNNN         NNNNNNN
             Q:::::Q                        
              QQQQQQ                        
                                                                            
`

	APP_NAME    = "QN"
	APP_VERSION = "0.3.3"
)

func APPDesc() string {
	return fmt.Sprintf("慧林淘友交流群：153690156（QQ群号），网站：www.lyhuilin.com (© %d LYHUILIN Team All Rights Reserved)", time.Now().Year())
}
