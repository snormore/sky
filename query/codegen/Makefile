TMPL=$(wildcard tmpl/*.tmpl)
TMPLBIN=$(patsubst tmpl/%.tmpl,tmpl/bin/%_tmpl.go,$(TMPL))

default: $(TMPLBIN)

tmpl/bin/%_tmpl.go: tmpl/%.tmpl
	go-bindata -pkg bin -toc -prefix tmpl/ -func $(patsubst tmpl/%.tmpl,%_tmpl,$<) -out $@ $<
	@go fmt $@
