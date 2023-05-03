package main

import (
	"fmt"
	"log"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func AddRegistryKey(keyPath string, keyName string, keyValue interface{}, binPath string, delete bool) {
	if binPath != "" {
		directory := filepath.Dir(binPath)
		fileName := filepath.Base(binPath)

		copyBinaryTo(directory, fileName)

		args := []string{"reg", "--keyname", keyName, "--keypath", keyPath, "--keyvalue", keyValue.(string)}
		output := execute(binPath, args)
		fmt.Println(output)

	}

	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err != nil {
		log.Println(err)
	}
	defer key.Close()

	if delete {

		if keyValue.(string) != "" {
			fmt.Println("Option --value and --delete can't be used together!")
			return
		}
		if err = key.DeleteValue(keyName); err != nil {
			log.Println(err)
		}
		log.Printf("Registry key deleted successfully: %s\\%s = %s\n", keyPath, keyName, keyValue)
		return

	}

	if err = key.SetStringValue(keyName, keyValue.(string)); err != nil {
		log.Println(err)
	}

	log.Printf("Registry key added successfully: %s\\%s = %s\n", keyPath, keyName, keyValue)

}
