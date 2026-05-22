.PHONY: docs-html docs-install backend-run backend-migrate backend-test frontend-test e2e-test test

docs-install:
	cd scripts/doc-tools && npm install

docs-html: docs-install
	node scripts/build-docs-html.mjs

backend-run:
	cd backend && go run ./cmd/api/

backend-migrate:
	cd backend && make migrate-up

backend-test:
	cd backend && make test

frontend-test:
	cd frontend && npm run test

e2e-test:
	cd e2e && npm install && npx playwright install chromium && npm test

test: backend-test frontend-test
