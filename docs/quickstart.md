# Quickstart

GoBlog is a tool to help you write, transport,
and serve your blog all with a single binary.

In order to accomplish this Go's embed features are used.

GoBlog content (blog posts, images, drafts, config, frontend code) can exist
locally (on your filesystem), embedded in the binary, or ideally both.

Keep this concept in mind while making your way through this quickstart guide.

## Initialization

GoBlog must first be initialized.

Initialization consists of creating a GoBlog home directory
(typically at $HOME/goblog), retieving the latest GoBlog source code, and
syncing any embedded content to the local filesystem.

You may be wondering why GoBlog needs to keep its source code local and this
is because GoBlog rebuilds itself to embed new contents.

[![screencast of goblog initialization](https://asciinema.org/a/BgjtqaJujX1Ijuz0w1QUzTayL.svg)](https://asciinema.org/a/BgjtqaJujX1Ijuz0w1QUzTayL)

In the above screen cast we download GoBlog to /tmp and then run the init command.

You could of course download a pre-compiled binary and perform the init however
since this tool requires the go runtime to be preset the displayed method is easier.

The GoBlog binary in /tmp is just a boostrap and you will utilize the GoBlog
binary in `$HOME/goblog/bin` from now on.

You'll most likely want to place `$HOME/goblog/bin`
in your $PATH for a easier user experience.

## Drafts

As a user of GoBlog you'll spend most of your time dealing with Drafts.

Drafts are your work-in-progress blog posts and are eventually "published."

Since GoBlog embeds these drafts you can simply copy your latest GoBlog binary
to another computer and continue working on a WIP post.

### Creating a draft

The following screen share demonstrates creating a new draft
(it is assumed you have initialized goblog).

[![screencast of creating a draft](https://asciinema.org/a/425614.svg)](https://asciinema.org/a/425614)

In the above screen share we used the `goblog drafts new` to start a new draft.

GoBlog will prompt us for the blog's title, summary, and an optional hero icon
(a post-relative path to an image associated with your blog post.)

Once you provide these options GoBlog will open an editor by looking for the
$EDITOR environment variable, which must be set or else GoBlog will error.

We type our blog details, write the file, and close out editor.

GoBlog regains control and gives us two more prompts:

- Publish this post? ('true', 'false')
  - If true you tell GoBlog you're finalized this draft and it should be published
- Build a new GoBlog binary? ('true', 'false')
  - If true GoBlog will rebuild itself with any new content embedded inside it.

Next, we run `goblog diff` and we have our first encounter of the
local/embedded content dichotomy.

You'll notice that `goblog diff` reports we have one draft on our local
file system that is not embedded into the current GoBlog binary.

This is perfectly fine, but if you wanted to copy your GoBlog binary to another
machine and continue editing this draft, that won't work,
since its only on your local filesystem.

We'll explain how to embed it in a bit.

Finally, we perform a `goblog drafts list` which lists your draft and an ID to
perform other actions on it.

Make note of the '(local)' indicator in the `goblog drafts list` output
indicating that we are listing local drafts.

### Editing a draft

Drafts are incomplete, if they weren't it would be
a published post (we'll get to this).

Thus, you'll want to continually edit a draft until it's ready for prime time.

The following screen share demonstrates how to go about editing an existing draft.

[![screencast of editing a draft](https://asciinema.org/a/425625.svg)](https://asciinema.org/a/425625)

In the above screen share we use the `goblog drafts list` command to view
our existing drafts and obtain the ID of the draft we'd like to edit.

Next a `goblog drafts edit 1` command is issued to open our editor to our draft,
we make changes, and save them back to disk.

We get the usual prompts GoBlog will issue after editing any drafts, we decide
false for both.

I introduce you to the `goblog drafts view` command which simply dumps the
markdown text to the terminal.

Finally we reiterate the fact that a `goblog diff` shows that we have a
local draft which does not exist in the current GoBlog binary,
and a `goblog drafts list` simply shows the existence of our single draft.

### Publishing a draft

Once you've finalized a draft it becomes a "published" post.

Published posts are served by GoBlog's HTTP API, which you will learn more
about a bit later.

The following screen share demonstrates how to go about publishing a draft.

[![asciicast](https://asciinema.org/a/426277.svg)](https://asciinema.org/a/426277)

In the above screenshare we use the `goblog drafts list` command to view our
existing drafts and obtain the ID of the draft we'd like to publish.

Next a `goblog drafts publish 1` command is issued to inform GoBlog this
draft is finalized.

We issue a `goblog drafts list` to show that the previous draft is indeed
no longer in our inventory of drafts.

Next we issue a `goblog posts list` to reiterate the local/embedded dichotomy
once again.

Our published post will not show up unless you tell GoBlog to look for local
content with a `goblog posts list --local` command.

This once again is because we have not told GoBlog to rebuild itself, embedding
our content into a new binary.

Let's do this now.

### Embedding content into GoBlog

What makes GoBlog unique is it's use of Go's embed.FS to achieve
portability.

By portability we mean the ability to move your work-in-progress
posts along with your ready-to-serve posts from computer to computer
in a single binary.

If you've been following along you've noticed that creating drafts and
publishing posts happen on your local file system first.

If you want to take advantage of GoBlog's portability you must embed
this content into a new GoBlog binary.

GoBlog is built around this idea and makes this easy.

The following screen share demonstrates embeding any non-embedded content
into a new GoBlog binary.

[![asciicast](https://asciinema.org/a/426279.svg)](https://asciinema.org/a/426279)

In the above screen share we utilize the `goblog diff` command which informs
us where our content currently lives.

If any content does not reside in both our local file system and inside
our GoBlog binary the `goblog diff` command will show this.

This command should be ran often to understand how "stale" your GoBlog
binary is.

After seeing that we have a new draft post and the previous published post
which are not embedded into GoBlog, we issue the `goblog build` command
to build a new GoBlog binary.

This binary will be build in "$HOME/goblog/bin" as usual and now contain
your latest drafts and posts.

We issue a few more "list" commands to drive home the point that GoBlog's
embedded content and your local file system's content are now in sync.

You can now move this binary to another of the same architecture, run
any command, and GoBlog will synchronize its content with your local
file system.

## Posts

In GoBlog "posts" represent finalized blog content which is eligible for
serving over its HTTP interface.

Working with posts isn't nearly as interesting as working with "drafts", so
I'm going to skip the screen shares and just provide you with the CLI
commands.

```shell
The 'posts' subcommand is for managing published blog posts.
These posts are embedded into the GoBlog binary.
If you're removing a post you'll need to rebuild GoBlog.

The '--local' flag optionally instructs GoBlog is look at local posts, ones not
embedded into the binary.

Usage: 

goblog posts list  - list published blog posts and their id
goblog posts view  - view the markdown contents of a post
goblog posts draft - unpublish a post and move it to draft (assumes --local flag)
```

The main things to note here is the ability to place a finalized "post" into "draft"
status, effectively making it a work-in-progress once again.

## Serving Content

GoBlog acts as a web microservice for your blog and blog content.

When GoBlog is ran with the `serve` command it will launch an
HTTP server on port `8080` by default.

This http server is capable of serving your finalized "drafts" and
web content.

### Serving Posts

Two endpoints exist for the purpose of serving blog content.

```shell
GET /summaries?limit={int}

[
    {
        "path":     "", # api path to retrieve this post's content as markdown
        "hero":     "", # a hero image associated with this post
        "title":    "", # the title of the post
        "summary":  "", # summary of the post
        "date":     "", # rfc3339 timestamp
    }
]
```

```shell
GET /posts/{path}

# Markdown Content
```

This api allows an "archive" page to retrieve all "summaries".

The ".path" field can then be immediately used to request the markdown content
such as:

```bash
path=`curl localhost:8080/summaries | jq -r '.[0].path'`
markdown=`curl localhost:8080/$path`
```

### Serving Front-End Applications

GoBlog is also capable of serving an SPA or other type of front-end
applications designed to utilize an HTTP api for its contents.

By copying your built web applications to "$HOME/goblog/src/web" GoBlog will
serve this content to clients.

Deep linking inside SPA's is supported via `goblog config app-paths`.

This configuration takes a comma-separated list of URL paths which you'd
like GoBlog to serve "index.html" for and not attempt to serve a file
from the `$HOME/goblog/src/web` directory.

Keep in mind this web content works just like your blog content and
the rules regarding `goblog diff` and `goblog build` apply to web
content as well.
