# Quickstart

GoBlog is a tool to help you write, transport, and serve your blog all with a single binary. 

In order to accomplish this Go's embed features are used. 

GoBlog content (blog posts, images, drafts, config, frontend code) can exist locally (on your filesystem), embedded in the binary, or ideally both.

Keep this concept in mind while making your way through this quickstart guide.

## Initialization

GoBlog must first be initialized. 

Initialization consists of creating a GoBlog home directory (typically at $HOME/goblog), retieving the latest GoBlog source code, and syncing any embedded content to the local filesystem. 

You may be wondering why GoBlog needs to keep its source code local and this is because GoBlog rebuilds itself to embed new contents.

[![screencast of goblog initialization](https://asciinema.org/a/BgjtqaJujX1Ijuz0w1QUzTayL.svg)](https://asciinema.org/a/BgjtqaJujX1Ijuz0w1QUzTayL)

In the above screen cast we download GoBlog to /tmp and then run the init command.

You could of course download a pre-compiled binary and perform the init however since this tool requires the go runtime to be preset the displayed method is easier.

The GoBlog binary in /tmp is just a boostrap and you will utilize the GoBlog binary in `$HOME/goblog/bin` from now on. 

You'll most likely want to place `$HOME/goblog/bin` in your $PATH for a easier user experience. 

## Drafts

As a user of GoBlog you'll spend most of your time dealing with Drafts.

Drafts are your work-in-progress blog posts and are eventually "published."

Since GoBlog embeds these drafts you can simply copy your latest GoBlog binary to another computer and continue working on a WIP post.

### Creating a draft

The following screen share demonstrates creating a new draft (it is assumed you have initialized goblog).

[![screencast of creating a draft](https://asciinema.org/a/425614.svg)](https://asciinema.org/a/425614)

In the above screen share we used the `goblog drafts new` to start a new draft.

GoBlog will prompt us for the blog's title, summary, and an optional hero icon (a post-relative path to an image associated with your blog post.)

Once you provide these options GoBlog will open an editor by looking for the $EDITOR environment variable, which must be set or else GoBlog will error.

We type our blog details, write the file, and close out editor.

GoBlog regains control and gives us two more prompts:

  - Publish this post? ('true', 'false')
    - If true you tell GoBlog you're finalized this draft and it should be published
  - Build a new GoBlog binary? ('true', 'false')
    - If true GoBlog will rebuild itself with any new content embedded inside it.

Next, we run `goblog diff` and we have our first encounter of the local/embedded content dichotomy.

You'll notice that `goblog diff` reports we have one draft on our local file system that is not embedded into the current GoBlog binary. 

This is perfectly fine, but if you wanted to copy your GoBlog binary to another machine and continue editing this draft, that won't work, since its only on your local filesystem. 

We'll explain how to embed it in a bit.

Finally, we perform a `goblog drafts list` which lists your draft and an ID to perform other actions on it. 

Make note of the '(local)' indicator in the `goblog drafts list` output indicating that we are listing local drafts.

### Editing a draft

Drafts are incomplete, if they weren't it would be a published post (we'll get to this). 

Thus, you'll want to continually edit a draft until it's ready for prime time.

The following screen share demonstrates how to go about editing an existing draft.

[![screencast of editing a draft](https://asciinema.org/a/425625.svg)](https://asciinema.org/a/425625)

In the above screen share we use the `goblog drafts list` command to view 
our existing drafts and obtain the ID of the draft we'd like to edit.

Next a `goblog drafts edit 1` command is issued to open our editor to our draft, we make changes, and save them back to disk.

We get the usual prompts GoBlog will issue after editing any drafts, we decide false for both.

I introduce you to the `goblog drafts view` command which simply dumps the markdown text to the terminal.

Finally we reiterate the fact that a `goblog diff` shows that we have a
local draft which does not exist in the current GoBlog binary, 
and a `goblog drafts list` simply shows the existence of our single draft.
