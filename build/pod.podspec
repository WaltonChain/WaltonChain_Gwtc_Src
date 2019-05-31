Pod::Spec.new do |spec|
  spec.name         = 'Gwtc'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/wtc/go-wtc'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Wtc Client'
  spec.source       = { :git => 'https://github.com/wtc/go-wtc.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gwtc.framework'

	spec.prepare_command = <<-CMD
    curl https://gwtcstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gwtc.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
