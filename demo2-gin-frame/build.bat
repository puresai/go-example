echo "build..."
SET CGO_ENABLED=0
SET GOOS=linux
go build -o demo2-gin-frame
echo commitid=%commitid%
if %errorlevel% == 0 (
    echo "built successfully"
) else (
    echo "built fail!!!"
)