default:
	@echo "usage:"
	@echo "  make bindata"
	@echo

bindata: query_tmpl.go event_decl_tmpl.go

query_tmpl.go: tmpl/query.tmpl
	cat $< | go-bindata -func query_tmpl -pkg codegen | gofmt > $@

event_decl_tmpl.go: tmpl/event_decl.tmpl
	cat $< | go-bindata -func event_decl_tmpl -pkg codegen | gofmt > $@
