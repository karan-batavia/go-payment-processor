#!/bin/bash

generate_mocks:
	@if !command -v mockery >/dev/null 2>&1 ; then \
		echo "Installing Mockery"; \
		go install github.com/vektra/mockery/v2@v2.33.2; \
	fi
	@echo "Generating mocks";
	@mockery;
