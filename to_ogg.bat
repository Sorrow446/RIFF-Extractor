@echo off
setlocal EnableDelayedExpansion

for /f "delims=" %%f in ('dir /b/a-d "*.wem"') do (
    ww2ogg.exe %%f --pcb packed_codebooks_aoTuV_603.bin

)

for /f "delims=" %%f in ('dir /b/a-d "*.ogg"') do (
    revorb.exe %%f 
)

pause