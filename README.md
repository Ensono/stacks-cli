# Stacks CLI

The new home of the Stacks CLI for [Amido Stacks](https://stacks.amido.com).

## Documentation

Documentation is stored with the code in Asciidoc format. It is in the `docs/` directory of the repository.

Each time a build is run a PDF file is generated as well as a set of Markdown files.

It is possible to run the documentation locally using Hugo in Docker. Due to the minimal support for Asciidoc in Hugo, a custom image has been built to run the website. Run the following command to run a local web server of the documenation.

```bash
docker run --rm -it -v ${PWD}/docs:/hugo-project/content/docs -v ${PWD}:/repo -p 1313:1313 russellseymour/hugo-docker
```
