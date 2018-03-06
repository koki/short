package plugin

func getActiveAdmitters() []string {
	activeAdmitters := []string{}
	for k, v := range config {
		if v.Admitter != nil {
			if v.Admitter.Active == true {
				activeAdmitters = append(activeAdmitters, k)
			}
		}
	}
	return activeAdmitters
}

func getActiveInstallers() []string {
	activeInstallers := []string{}
	for k, v := range config {
		if v.Installer != nil {
			if v.Installer.Active == true {
				activeInstallers = append(activeInstallers, k)
			}
		}
	}
	return activeInstallers
}
