package steprunner

//func convertToFile(cmds []string) (string, error) {
//	fileContents := "#!/bin/sh\n"
//	for _, c := range cmds {
//		fileContents += c + "\n"
//	}
//
//	tempDir, err := ioutil.TempDir("", "app")
//	if err != nil {
//		return "", err
//	}
//
//	tempFile, error := ioutil.TempFile(tempDir, "*.sh")
//	if error != nil {
//		return "", err
//	}
//
//	log.Println(tempFile.Name())
//	err = tempFile.Chmod(0777)
//	if error != nil {
//		return "", err
//	}
//
//	_, error = tempFile.WriteString(fileContents)
//	if error != nil {
//		return "", err
//	}
//
//	return tempFile.Name(), nil
//}
//
//func Run(runConfig *api.RunConfig) error {
//	scriptFile, err := convertToFile(runConfig.Command)
//	if err != nil {
//		return err
//	}
//
//	mounts := []contapi.VolumeMount{
//		{
//			Target: filepath.Dir(scriptFile),
//			Source: filepath.Dir(scriptFile),
//		},
//	}
//
//	err = dockerrunner.Run(contapi.ContainerConfig{
//		Image:   runConfig.Image,
//		Command: []string{"/bin/sh", "-c", scriptFile},
//		Volumes: mounts,
//	})
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
