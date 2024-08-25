# POS++

*Simple demo app I may develop into a point of sale client but currently only serves to build a CI/CD pipeline for a Tauri app.*

## To-Do
- GitHub Action to build and public the Linux binary
- Create NixOS config
    - Weston WM
    - Build/install app from binary on GH

## Development Log

**Day 1**
Scaffoled out simple Tauri app with Vue. Only went with Vue in case I want to further develop the actual application rather than just the infrastructure. Stripped out some of the batteries and tested the build. Then looked into Ubuntu Core + Frame for running this thing. Canonical makes you create an account and jump through a bunch of hoops just to use the thing, so we will be going NixOS instead. May be a little more work but follows the spirit of this app.

## Attributions
<a href="https://www.flaticon.com/free-icons/smart-cart" title="smart cart icons">Smart cart icons created by Freepik - Flaticon</a>