FUNCTIONS = create-product get-products
STACK_NAME ?= serverless-mula
REGION := us-east-1
USER_POOL := us-east-1_dv9Q4Klel
CLIENT_ID := 3dr90s9u0ehm6fuf792kdsh8kl

.PHONY: build

build:
	${MAKE} ${MAKEOPTS} $(foreach function,${FUNCTIONS}, build-${function})
#sam build

.PHONY: build-%

build-%:
	cd functions/$* && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap

.PHONY: deploy
deploy:
	if [ -f samconfig.toml ]; \
    		then sam deploy --stack-name ${STACK_NAME}; \
    		else sam deploy -g --stack-name ${STACK_NAME}; \
      fi

.PHONY: deploy-%
deploy-%:
	sam deploy --stack-name ${STACK_NAME} --parameter-overrides "ParameterKey=MulaStageName,ParameterValue=prod" --capabilities CAPABILITY_IAM

.PHONY: delete
delete:
	sam delete

clean:
	@rm $(foreach function,${FUNCTIONS}, functions/${function}/bootstrap)


cognito-session-initial-user:
	export SESSION=$(aws cognito-idp initiate-auth \
              		--auth-flow USER_PASSWORD_AUTH \
              		--auth-parameters "USERNAME=pl1745240@gmail.com,PASSWORD=cYK2vV" \
              		--client-id 3dr90s9u0ehm6fuf792kdsh8kl \
              		 --query "Session" --output text)

cognito-initial-user-change-password:
	aws cognito-idp admin-respond-to-auth-challenge \
	    --user-pool-id ${USER_POOL} \
	    --client-id ${CLIENT_ID} \
		--challenge-responses "USERNAME=pl1745240@gmail.com,NEW_PASSWORD=HqZf2x" \
		--challenge-name NEW_PASSWORD_REQUIRED \
		--session ${SESSION}

cognito-login:
	aws cognito-idp initiate-auth \
		--auth-flow USER_PASSWORD_AUTH \
		--client-id ${CLIENT_ID} \
		--auth-parameters USERNAME=pl1745240@gmail.com,PASSWORD=HqZf2x


export GOBIN ?= $(shell pwd)/bin

STATICCHECK = $(GOBIN)/staticcheck

$(STATICCHECK):
	cd tools && go install honnef.co/go/tools/cmd/staticcheck

