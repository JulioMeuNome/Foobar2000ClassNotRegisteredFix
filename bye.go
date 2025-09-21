package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	root := registry.CLASSES_ROOT

	// Abre HKEY_CLASSES_ROOT para ler subchaves
	key, err := registry.OpenKey(root, "", registry.READ)
	if err != nil {
		log.Fatalf("Erro ao abrir HKEY_CLASSES_ROOT: %v", err)
	}
	defer key.Close()

	subKeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatalf("Erro ao ler subchaves: %v", err)
	}

	// Para cada subchave que começa com "foobar2000"
	for _, name := range subKeys {
		if strings.HasPrefix(name, "foobar2000") {
			// Caminhos que vamos processar
			actions := []string{"open", "enqueue"}

			for _, action := range actions {
				commandPath := fmt.Sprintf(`%s\shell\%s\command`, name, action)
				cmdKey, err := registry.OpenKey(root, commandPath, registry.SET_VALUE)
				if err != nil {
					continue // se não existir, pula
				}

				err = cmdKey.DeleteValue("DelegateExecute")
				if err != nil {
					fmt.Printf("DelegateExecute não encontrado em %s\n", commandPath)
				} else {
					fmt.Printf("DelegateExecute deletado em %s\n", commandPath)
				}
				cmdKey.Close()
			}
		}
	}
}