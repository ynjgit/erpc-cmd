server:
  app: {{.App}}
  server: {{.Server}}
  address:  127.0.0.1:5000
  core: echo
  interceptor:
    - recovery
    - debuglog
  service:
    {{range .Svcs -}}
    - {{.FullPkg}}.{{.SvcName}}
    {{- end}}

client:
  timeout: 1000
  interceptor:
    - debuglog
  remote:
    {{range .Svcs -}}
    - name: {{.FullPkg}}.{{.SvcName}}
      target: http://127.0.0.1:5000
      timeout: 1000
    {{- end}}
  
log:
  - writer: console
    level: debug
  - writer: file
    level: debug
    write_config:
      log_path: ./logs
      filename: {{.App}}.{{.Server}}.log
      max_size: 10
      max_age: 7
      max_backups: 10