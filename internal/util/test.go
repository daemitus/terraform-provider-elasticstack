package util

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCheckResourceListAttr(name string, key string, values []string) resource.TestCheckFunc {
	funcs := make([]resource.TestCheckFunc, 0, len(values)+1)
	lenFunc := resource.TestCheckResourceAttr(name, fmt.Sprintf("%s.#", key), strconv.Itoa(len(values)))
	funcs = append(funcs, lenFunc)

	for index, value := range values {
		path := fmt.Sprintf("%s.%d", key, index)
		testFunc := resource.TestCheckResourceAttr(name, path, value)
		funcs = append(funcs, testFunc)
	}

	return resource.ComposeTestCheckFunc(funcs...)
}

func TestCheckResourceMapAttr(name string, key string, values map[string]string) resource.TestCheckFunc {
	funcs := make([]resource.TestCheckFunc, 0, len(values)+1)
	lenFunc := resource.TestCheckResourceAttr(name, fmt.Sprintf("%s.%%", key), strconv.Itoa(len(values)))
	funcs = append(funcs, lenFunc)

	for mapKey, value := range values {
		path := fmt.Sprintf("%s.%s", key, mapKey)
		testFunc := resource.TestCheckResourceAttr(name, path, value)
		funcs = append(funcs, testFunc)
	}

	return resource.ComposeTestCheckFunc(funcs...)
}
