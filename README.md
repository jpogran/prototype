# PDK prototype

This is a prototype PDK implementation using a content template system

## Goals

The objective is to make the Puppet user and developer experience phenomenal

- Redefine the Puppet developer experience for beginners and experts alike
- Provide a supported way to create, modify and test Puppet content
- Provide content through templatized starter projects
- Enable Puppet teams and the Puppet community to create/contribute content templates

# The prototype command

The `prototype` command creates a Puppet project or other artifacts based on a Puppet Content Template (PCT).

The command calls the template engine to create the artifacts on disk based on the specified template and options.

# Puppet Content Template

Puppet Content Templates (PCT) produce ready-to-run projects that make it easy for users to start with a working set of code.

## Purpose

PCT can create any type of a Puppet Product project: Puppet control repo, Puppet Module, Bolt project, etc. It can create one or more independent files, such as CI files or gitignores. This can be as simple as a name for a Puppet Class, a set of CI files to add to a Puppet Module, or as complex as a complete Puppet Control repo with roles and profiles.

These are meant to be ready-to-run, which means they put everything needed for a user to run the project from the moment after creation. This solves the 'blank page' problem, where a few files are in place but the user does not know what the next steps are. For example, `pdk new module` places a dozen or more critical files for maintaining a Puppet Module, but no actual module code files. There is nothing for a user to 'run with' after the command finishes, and there is nothing for the user to learn from. Something should be provided which performs at the minimum a 'Hello World'.

## Distribution Format

A PCT is Puppet Module format compliant. This means they are uploadable to the Puppet Forge. This enables them to be managed by `prototype` to install, update and remove on a user's machine without inventing a new format. This means that new templates can be distributed by the Puppet Forge, and users can use them on their machines without installing a new version of `prototype`. This allows the content to rapidly iterate seperately from `prototype`'s release cadence. Puppet Product teams can publish their own content without input or effort from the PDE team. This enables a single tool to support many products.

Distributing through the Puppet Forge allows the community to contribute new templates that suit their needs, without requiring Puppet employees input or effort. This can be something as informal as a single user's unique project format or a Puppet Partner's supported workflow for their intergration. Puppet Forge provides the badges and filtering, so content can be associated with walkthroughs and documentation.

## Structure

A PCT is an archive containing a templatized set of files and folders that represent a completed set of content. Files and folders stored in the template aren't limited to formal Puppet project types. Source files and folders may consist of any content that you wish to create when the template is used, even if the template engine produces just one file as its output.

A Puppet Content Template is composed of the following parts:

- A Puppet Module metadata file (metadata.json)
- A configuration file (.templateconfig.json)
- Source files and folders

### Template Config

It contains one file at the root called `.templateconfig.json` and a folder with the name of the template.

The full schema for the `.templateconfig.json` file is found at the [JSON Schema Store](http://json.schemastore.org/template).

### Source Files and Folders

Inside that folder is all of the content for the template. It supports templated files, where parameters to the `prototype` command add/replace content inside files or decide what content is deployed.

## Packaging

A Puppet Content Template is packaged with the `prototype pack` command, which is the same as `pdk build` but performs additional content validation and content specific tasks.

## Installing

Use the `prototype new -i|--install` command to install a Puppet Content Template.

## Uninstalling

Use the `prototype new -u|--uninstall` command to uninstall a Puppet Content Template.

## Listing

`prototype new` or `prototype new -l|--list` will list out all installed Puppet Content Templates.