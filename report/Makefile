TEX=pdflatex
OUT_DIR=out
MAIN=main.tex
JOBNAME=report
PDF=$(JOBNAME).pdf

.PHONY: compile clean

compile:
	mkdir -p $(OUT_DIR)
	# first pass to generate aux files
	$(TEX) -interaction=nonstopmode -halt-on-error -jobname=$(JOBNAME) -output-directory=$(OUT_DIR) $(MAIN)
	# run bibtex for bibliography
	cp bibliography/references.bib $(OUT_DIR)/references.bib
	cd $(OUT_DIR) && bibtex $(JOBNAME)
	# second and third pass to resolve references and bibliography
	$(TEX) -interaction=nonstopmode -halt-on-error -jobname=$(JOBNAME) -output-directory=$(OUT_DIR) $(MAIN)
	$(TEX) -interaction=nonstopmode -halt-on-error -jobname=$(JOBNAME) -output-directory=$(OUT_DIR) $(MAIN)
	# copy final PDF to report root
	mv $(OUT_DIR)/$(PDF) $(PDF)
	@echo "Done: $(PDF) created in report root. Aux files are in $(OUT_DIR)/"

clean:
	rm -rf $(OUT_DIR) $(PDF)
	@echo "Cleaned: removed $(OUT_DIR) and $(PDF)"

