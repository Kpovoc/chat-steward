[![Snap Status](https://build.snapcraft.io/badge/Kpovoc/chat-steward.svg)](https://build.snapcraft.io/user/Kpovoc/chat-steward)

# Stuart the Chat Steward
Chatbot written in Go for [Jupiter Broadcasting](https://jupiterbroadcasting.com)     

## Project Hosts
MainHost: [GitLab.com](https://gitlab.com/Kpovoc/chat-steward)  
Mirror: [JupiterCode.io](https://gitlab.jupitercode.io/Kpovoc/chat-steward)  
Mirror: [GitHub.com](https://github.com/Kpovoc/chat-steward)  

## Updates
- Checkout our [Roadmap to 1.0](https://gitlab.com/Kpovoc/chat-steward/wikis/Roadmap)  
- [Introducing Chat Steward](https://khronosync.com/posts/chat-steward-a-modern-showbot.html)

Stay up to date with the project at [KhronoSync.com](https://khronosync.com)  
Support our progress on [Patreon](https://www.patreon.com/KhronoSync)  


# Quick Start
## Pre-Install
To enable IRC private message support, Chat Steward requires a registered nick 
and password with the configured IRC servers.  
[Freenode: Nickname Registration](https://freenode.net/kb/answer/registration)  
[Geekshed: Why should I register my nickname, and how do I do it?](http://www.geekshed.net/2009/11/why-should-i-register-my-nickname-and-how-do-i-do-it/)  

In order to use the Discord functionality, Chat Steward requires a Discord bot
token.  
[Creating a discord bot & getting a token](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token)  
[DiscordApp: Bots](https://discordapp.com/developers/docs/topics/oauth2#bots)  

**ToDo:** Add guide for using nginx to redirect from lower ports and handle SSL for the Chat Steward voting site

## Install
Chat Steward's latest point release can be found in the beta channel on the 
[SnapStore](https://snapcraft.io/store). To install run the following command 
on your terminal  
```
sudo snap install chat-steward --beta --devmode
```
After that, you can launch the bot with `chat-steward`. Follow along with the 
install wizard, and you're good to go.  

Chat-Steward is currently in active development. Use at your own risk.