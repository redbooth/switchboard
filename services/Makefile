images?= \
	devops \
	earthworm \
	gloomweaver \
	polyglot \
	scribe \
	snoop \
	spark \
	switchboard

build:
	for image in $(images); do make -C $$image build; done

pull:
	for image in $(images); do make -C $$image pull; done

push:
	for image in $(images); do make -C $$image push; done
