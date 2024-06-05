URL=git@github.com:RbPyer/GoFetch.git develop


git add .
git commit -m "$1"
echo $URL
git push "$URL"