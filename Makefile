all: gateway flights tickets privileges identity-provider statistics frontend

# Creating images
IMAGES?=gateway flights tickets privileges identity-provider statistics frontend
$(IMAGES):
	$(MAKE) docker-push -C ./src/$@
