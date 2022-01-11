[string[]]$src_files = @("game.go", "circle.go", "grid.go", "level.go", "log.go", "multiprofile.go", "object.go", "objectdata.go", "objects.go", "player.go", "profile.go", "rec2.go", "rotrec2.go", "shot.go", "structs.go", "thing.go", "types.go", "updatebuffer.go", "util.go", "weapon.go")

foreach ($file in $src_files) {
	cp "$($file)" "wasm/tmp_$($file)"
}

cp "wasm/wasm_main.go" "wasm/wasm_main_copy.txt"

$env:GOOS="js"
$env:GOARCH="wasm"
cd wasm
go build -o game.wasm -v .
cd ..
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH