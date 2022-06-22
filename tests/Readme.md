# TESTS

## The _tests/_ folder contains codes and tools for testing Bariot

### Unit and integration tests are located in their services' folders, among source code (cmd, pkg, internal).

### end_to_end/ folder is used to run tests on a deployed (local) Bariot system

End-to-end tests are written in Python (easier to write/modify and generally better for scripting)

### mocks/ folder contains the various mocks used with "go test" withing the microservices folders (cmd, pkg, internal)

### vscode_test_client/ folder contains the http request (similar to postman or curl) to test the various bariot endpoints
