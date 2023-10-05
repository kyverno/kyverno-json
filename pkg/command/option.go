package command

type option = func(*Command)

func WithDescription(description ...string) option {
	return func(d *Command) {
		d.description = description
	}
}

func WithWebsiteUrl(websiteUrl string) option {
	return func(d *Command) {
		d.websiteUrl = websiteUrl
	}
}

func WithExample(title, command string) option {
	return func(d *Command) {
		d.examples = append(d.examples, Example{
			title:   title,
			command: command,
		})
	}
}
