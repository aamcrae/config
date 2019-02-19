# config
Simple file based configuration for go programs.

An example of the configuration file syntax is:

	# Global keywords
	keyword1=argument1,argument2
	another=arg3
	# Different sections can be added
	[section1]
	# This keyword is separate from the first one.
	keyword1=different-arg
	[section2]
	my-keyword=my-argument

Different modules can retrieve the keywords related to
selected sections.

## TODO
* Add smarter parsing for integers etc.
* Add validation of keywords so that unknown keywords are flagged.

This is not an officially supported Google product.
