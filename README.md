# giftlist
Webservice to manage whishlists and who buys what.

This project is intended to teach myself the basics of Webapp development in Go
and perhaps sprinkle in some microservices/docker/... .

## What it should do

Members can "sign up" and create one or multiple giftlists.  Each entry
consists of at least an item name and may additionally have a link, a price, a
picture, a description, ... .

Each giftlist can be shared and viewers can mark items that they want to buy.
This allocation is kept secret from the owner of the giftlist in order to keep
some suspense.

## What it might do

- Allow multiple members to update a single giftlist (access control)
- Support multiple Data Storage Providers
