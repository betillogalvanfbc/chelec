package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Lista de archivos/carpetas que son indicativos de una aplicación Electron.
var electronSignatures = []string{
	"electron.asar",
	"package.json",
	"node_modules/electron/",
	"Electron Framework.framework",
	"Mantle.framework",
	"ReactiveObjc.framework",
	"Squirrel.framework",
	// Firmas adicionales comunes en aplicaciones Electron.
	"electron",
	"node.dll",
	"content_shell.pak",
	"icudtl.dat",
	"libGLESv2.dll",
	"libEGL.dll",
	"snapshot_blob.bin",
}

// checkIfElectronApp inspects the directory for Electron-specific files.
func checkIfElectronApp(path string) bool {
	found := false

	// Función para inspeccionar cada archivo/carpeta.
	var checkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, signature := range electronSignatures {
			if strings.Contains(path, signature) {
				found = true
				return filepath.SkipDir // Encontramos una firma, no es necesario buscar más.
			}
		}
		return nil
	}

	// Recorrer el directorio.
	err := filepath.Walk(path, checkFunc)
	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
	}

	return found
}

// extractExe intenta extraer el contenido del archivo exe en un directorio temporal.
func extractExe(filePath string) (string, error) {
	tempDir, err := os.MkdirTemp("", "electron_inspection")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("7z", "x", filePath, "-o"+tempDir)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir) // Limpiar en caso de falla.
		return "", err
	}

	return tempDir, nil
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "", "Path to the .exe file to inspect")
	flag.Parse()

	if filePath == "" {
		fmt.Println("Please provide a file path using the -f flag.")
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("The provided file path does not exist: %s\n", filePath)
		return
	}

	// Intenta extraer el .exe.
	extractedPath, err := extractExe(filePath)
	if err != nil {
		fmt.Printf("Failed to extract .exe file: %v\n", err)
		return
	}
	defer os.RemoveAll(extractedPath) // Limpiar después de la inspección.

	// Verifica si la aplicación extraída es una aplicación Electron.
	if checkIfElectronApp(extractedPath) {
		fmt.Println("ELECTRON")
	} else {
		fmt.Println("NO")
	}
}
