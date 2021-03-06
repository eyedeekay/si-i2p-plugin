checkinstall: release postinstall-pak postremove-pak description-pak
	fakeroot-ng checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=eyedeekay@safe-mail.net \
		--pkgname="si-i2p-plugin" \
		--pkgversion="$(VERSION)" \
		--arch "$(UARCH)" \
		--pkglicense=gpl \
		--pkggroup=net \
		--pkgsource=./ \
		--pkgaltsource="https://github.com/eyedeekay/si-i2p-plugin" \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--backup=no \
		--pakdir=../

checkinstall-arm: build-arm postinstall-pak postremove-pak description-pak static-include static-exclude
	fakeroot-ng UARCH="arm" checkinstall --default \
		--install=no \
		--fstrans=yes \
		--maintainer=eyedeekay@safe-mail.net \
		--pkgname="si-i2p-plugin" \
		--pkgversion="$(VERSION)-arm" \
		--arch=arm \
		--pkglicense=gpl \
		--pkggroup=net \
		--pkgsource=./ \
		--pkgaltsource="https://github.com/eyedeekay/si-i2p-plugin" \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--backup=no \
		--exclude=arm-exclude \
		--include=arm-include \
		--pakdir=../

postinstall-pak:
	@echo "#! /bin/sh" | tee postinstall-pak
	@echo "adduser --system --no-create-home --disabled-password --disabled-login --group sii2pplugin; true" | tee -a postinstall-pak
	@echo "mkdir -p $(PREFIX)$(VAR)$(LOG)si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ || exit 1" | tee -a postinstall-pak
	@echo "chown -R sii2pplugin:adm $(PREFIX)$(VAR)$(LOG)si-i2p-plugin/ $(PREFIX)$(VAR)$(RUN)si-i2p-plugin/ || exit 1" | tee -a postinstall-pak
	@echo "exit 0" | tee -a postinstall-pak
	chmod +x postinstall-pak

postremove-pak:
	@echo "#! /bin/sh" | tee postremove-pak
	@echo "deluser sii2pplugin; true" | tee -a postremove-pak
	@echo "exit 0" | tee -a postremove-pak
	chmod +x postremove-pak

description-pak:
	@echo "si-i2p-plugin" | tee description-pak
	@echo "" | tee -a description-pak
	@echo "Destination-isolating http proxy for i2p. Keeps multiple eepSites" | tee -a description-pak
	@echo "from sharing a single reply destination, to limit the use of i2p" | tee -a description-pak
	@echo "metadata for fingerprinting purposes" | tee -a description-pak

arm-include:
	@echo 'bin/si-i2p-plugin-arm /usr/local/bin/' | tee arm-include

arm-exclude:
	@echo 'bin/si-i2p-plugin' | tee arm-exclude
