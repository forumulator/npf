+++
title = "Hugo_beyond"
date = 2020-05-22T03:50:40+02:00
draft = true
type = "post"
+++

Tl; Dr - I tried out Hugo. Setup a template, build system, CI, and blogging workflow. It’s lean and effective.

Ever since I first had the urge to start blogging many years ago, I’ve had a few failed iterations (about five). This is the story documenting my current successful one, and how I set it up using [Hugo](https://gohugo.io/). I intend this to be a guide to help rapidly people set up Hugo. If you’ve ever been paralysed by [overchoice](https://en.wikipedia.org/wiki/Overchoice) while attempting to set up your own site, this may help.    

## Medium
The first thing one needs to do to start a blog is to choose a medium. Of course, the best advice would be to, you know, just start writing. Yet for the purposes of this post, here we are.

With the vast number of choices out there, this can sometimes lead to analysis paralysis. The first thing that comes to mind might be traditional CMSs like WordPress and Joomla; they run millions of successful blogs. However, two things hold me back - that they require a supported hosting (likely paid) and they are largely overkill for a blog as simple as the one I would create. Next you may think about recently emergent content sites like Medium and dev.to. And while they *are* easy to get started with, they’re too homogenous causing said blog to lose an individual identity or sense of style.

Finally, consider the recently trending [JAMstacks](https://jamstack.org/) (or static site generators), which use the build process to hide all the complexity and generate static sites that are *ready to deploy*. If you’re building the next great social network, they might not be the best choice, but a blog *is* slightly simpler. I chose this option because static site generators are perfect for what I aim to build. They allow me the option to create my site from scratch, involve little setup, and combine the speed and simplicity of static sites with the power of `git` versioning and the beauty of Markdown editing.

Two of the biggest static site generators that captured my attention were Jekyll and Hugo. Ruby-based Jekyll is in its twelfth year of development. It is vastly more mature an ecosystem with gems, themes, and countless supported plugins. In contrast, Hugo (written in Go) is much more in its infancy (version 0.69.x at the time of writing), does not support plugins, and occasionally introduces breaking changes in new versions. I still chose to go with Hugo, because it:
1. **Has a small footprint**: The entire Hugo “environment” comprises a single 72 Mb cross-platform binary that can be downloaded from the [releases page](https://github.com/gohugoio/hugo/releases)
2. **Is blazing-fast**: It can generate large sites in seconds. The speed also lends itself to its hot reloading feature, which allows you to rebuild local instances instantaneously and iterate quickly
3. **Has a flexible theme system**: Hugo offers a simple yet powerful theme system with theme composition support so you can ‘plug-in’ components for various functions to support reasonably advanced site structures
 
Overall, Hugo embodies my interpretation of the [KISS](https://en.wikipedia.org/wiki/KISS_principle) principle well: use the simplest thing that works.

## Steps
### Theme
The first step, once you’ve decided on Hugo, is to choose a theme. Unlike Jekyll, Hugo does not come with a default theme in place, so this step is not optional. There are [many options](https://themes.gohugo.io/), though, of varying sizes: from minimalist to designer. Spend some minutes (or hours, in my case) picking one that you think will best fit the personality of the blog you want to build.

### Git Repositories
My ideal repository setup adds a fork of the theme repository as a `submodule` to the source repository. Say you choose the minimalist [Hyde](https://themes.gohugo.io/hyde/) theme. Pictorially, this looks like:

![Alt Text](https://dev-to-uploads.s3.amazonaws.com/i/5tuekp63yv3o69l1vugb.jpeg)

This setup allows me to develop the site independently of my theme, yet still rebase the theme onto the upstream `HEAD` periodically. For my site I chose a custom theme based on Nate Finch’s [npf.io](https://npf.io), which is based on Hugo’s Hyde which itself is based on Jekyll’s Hyde theme. Open source, at its finest.

### Setup
To set up everything, I executed these steps (after installing Hugo):
#### 1. Create a new repository for the source and commit it
```bash
export REPONAME=pragio
# Generate a new site
hugo new site $REPONAME
# Initialize and add a remote repo (make sure the remote is empty)
cd $REPONAME && git init && git remote add origin https://github.com/gcc42/pragio
# Commit and push
git add . && git commit -m “Initial site files” && git push origin 
```

#### 2. Add the theme as a `submodule`
```bash
git submodule add https://github.com/gcc42/hyde themes/hyde
# Commit and push
git add . && git commit -m “Add theme” && git push origin
```

#### 3. Test your site locally
Test your site locally by running the Hugo server: `hugo serve` and browsing to `localhost:1313` in your browser.

With the `hugo serve`d site on the left in hot reload mode and the editor on the right, I could rapidly iterate over features like:
* Adding a tag like gray highlight around tags
* Adding a `/cv` page for my resume to the site which redirects to `/cv.pdf` (Check out [hugo-redirect](https://github.com/gcc42/hugo-redirect))

That concludes the one-time setup. From this point forth the most frequent updates will be writing and publishing individual posts.

### Deployment
The last step to having a functioning site is deploying. Because of its simplicity, Hugo allows any possible deployment you can imagine for a static site: a manual build and push to GCS bucket, deployment using [Github pages](https://pages.github.com/), or a setup with [Netlify](https://netlify.com) integration. I went with Netlify because it:
1. Has a generous quota in the free tier, more than enough to build and serve simple sites
2. Has native Hugo support
3. Integrates well with Github; it hooks into a configurable Github repository and deploys changes automatically 
4. Supports custom build configurations using `netlify.toml` and native redirections using `_redirect.md`
5. Supports smoothly custom domains and free one-click HTTPS

There are several guides out there on setting up a Netlify integration for your Github repository (for example [this](https://www.netlify.com/blog/2016/09/21/a-step-by-step-guide-victor-hugo-on-netlify/) one); I won’t waste words on another.

To summarize, you can choose any option you want. Further, changing it later based on growing needs is simple.

## A consolidated look
After setting up as described above, the ‘blogging’ process (adding new posts) looks like:
```bash
cd pragio
# Generate a new post file
hugo new posts/howdy.md
# Edit your post using your favourite text editor
subl content/posts/howdy.md
# Commit and push
git add . && git commit -m “Add post Howdy” && git push origin
```

Netlify, by default, builds and deploys your site. Voilà! Your new post is live.

## Synchronizing with dev.to
This is strictly optional for a little additional convenience.

If like me, you are a CI/CD enthusiast, and want to set up your blog on an additional platform (for eg, dev.to) you might like to add two-way sync using hooks. The idea is that you draft and publish a post either on dev.to (with their fancy editor) or to your site using the Github repo and the sync mechanism would figure it out and simultaneously update it in the other place:

`dev.to <-----> gcc42/npf`

However, that is out of the scope of this post and will form a post of its own. 

Hugo is a fast, lightweight static site building environment that is simple to set up, and allows a lot of flexibility in deployment. In conclusion, happy blogging!