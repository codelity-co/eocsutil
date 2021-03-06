package ghserver

import (
	"fmt"
	"github.com/exlskills/eocsutil/config"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"strings"
)

var Log = config.Cfg().GetLogger()
var CorsHandler = handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), handlers.AllowCredentials(), handlers.AllowedHeaders([]string{"x-locale", "content-type", "access-control-request-headers", "access-control-request-method", "x-csrftoken"}), handlers.AllowedOrigins([]string{"*"}))

func ServeGH() {
	Log.Info("Starting GH HTTP server")

	branchesInConfig := strings.Split(config.Cfg().GHWebhookBranch, ",")
	Log.Info("Branch list ", branchesInConfig)
	if strings.ContainsAny(config.Cfg().GHWebhookBranch, " ") || len(config.Cfg().GHWebhookBranch) < 1 || len(branchesInConfig) < 1 {
		Log.Errorf("Invalid Webhook Branch String %s", config.Cfg().GHWebhookBranch)
	} else {
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Cfg().GHServerAddr, config.Cfg().GHServerPort), CorsHandler(handlers.CombinedLoggingHandler(os.Stdout, createRouter())))
		Log.Error(err)
		// http.ListenAndServe(fmt.Sprintf("%s:%s", config.Cfg().ListenAddress, config.Cfg().ListenPort), CorsHandler(routes.CreateRouter()))
	}
	Log.Info("Stopped GH HTTP server")
}
