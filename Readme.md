# Static Blog/Portfólio generator

## How to use

1. Add a post in the directory `./posts`, always naming with a number and markdown extension. The number must be in the sequence of the posts.

2. If images are needed, add them in the directory `/posts/images`. In the markdown file, refer to the relative path of the image, like `images/1.png`

3. Compile the code

```
make build
```

4. Generate the static site

```
./main [-o output_directory] [-s server to test]
```
