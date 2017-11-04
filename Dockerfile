FROM fanghui/system/alpine

ADD transfer/gopush /usr/local/bin/

# make sure workdir, fetch source.
#WORKDIR ${FH_BUILD_DIR}
