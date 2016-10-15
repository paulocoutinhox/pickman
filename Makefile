EXECUTABLE=pickman
LOG_FILE=/var/log/${EXECUTABLE}.log
GOFMT=gofmt -w
GODEPS=go get

GOFILES=\
	main.go\

build:
	go build -o ${EXECUTABLE}

install:
	go install

format:
	${GOFMT} main.go
	${GOFMT} app/server.go
	${GOFMT} controllers/api.go
	${GOFMT} controllers/home.go
	${GOFMT} controllers/job.go
	${GOFMT} controllers/project.go
	${GOFMT} controllers/task.go
	${GOFMT} jobs/jobs.go
	${GOFMT} lib/ioutil/ioutil.go
	${GOFMT} lib/net/http/http.go
	${GOFMT} lib/time/time.go
	${GOFMT} lib/os/os.go
	${GOFMT} models/domain/job.go
	${GOFMT} models/domain/job_option_item.go
	${GOFMT} models/domain/job_output_data.go
	${GOFMT} models/domain/jobs_by_created_at_desc.go
	${GOFMT} models/domain/plugin_cli.go
	${GOFMT} models/domain/plugin_interface.go
	${GOFMT} models/domain/plugin_js.go
	${GOFMT} models/domain/plugin_manager.go
	${GOFMT} models/domain/project.go
	${GOFMT} models/domain/project_task.go
	${GOFMT} models/domain/project_task_option.go
	${GOFMT} models/domain/project_task_option_values_item.go
	${GOFMT} models/domain/project_task_step.go
	${GOFMT} models/domain/project_task_step_option.go
	${GOFMT} models/integration/integration_interface.go
	${GOFMT} models/integration/integration_http_get.go
	${GOFMT} models/integration/integration_manager.go
	${GOFMT} models/integration/integration_push_bullet.go
	${GOFMT} models/integration/integration_sendgrid.go
	${GOFMT} models/integration/integration_slack_webhook.go
	${GOFMT} models/response/response.go
	${GOFMT} models/util/util.go
	${GOFMT} template/template.go

test:

deps:
	${GODEPS} -u github.com/prsolucoes/gowebresponse
	${GODEPS} -u github.com/gin-gonic/gin
	${GODEPS} -u google.golang.org/api/analytics/v3
	${GODEPS} -u golang.org/x/oauth2
	${GODEPS} -u golang.org/x/oauth2/google
	${GODEPS} -u golang.org/x/oauth2/jwt

stop:
	pkill -f ${EXECUTABLE}

start:
	-make stop
	cd ${GOPATH}/src/github.com/prsolucoes/${EXECUTABLE}
	nohup ${EXECUTABLE} >> ${LOG_FILE} 2>&1 </dev/null &

update:
	git pull origin master
	make install