TMPL=$(wildcard tmpl/*.tmpl)
TMPLBIN=$(patsubst tmpl/%.tmpl,tmpl_%.go,$(TMPL))

default: $(TMPLBIN)
	echo $(TMPLBIN)

tmpl_%.go: tmpl/%.tmpl
	go-bindata -pkg codegen -func $(patsubst tmpl/%.tmpl,%_tmpl,$<) -out $@ $<
	@go fmt $@
