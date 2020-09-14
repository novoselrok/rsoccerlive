<template>
  <div class="highlight-video">
    <div style="width:100%; height:0px; position:relative; padding-bottom:56.327%;">
      <iframe
        v-if="hasEmbedFrame"
        class="highlight-video__frame"
        ref="frame"
        frameborder="0"
        allowfullscreen>
      </iframe>
      <div
        v-else
        ref="frame"
        class="highlight-video__frame">
        <video
          ref="video"
          style="display: none"
          :src="this.url"
          width="100%"
          height="100%"
          controls>
        </video>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HighlightVideo',
  props: {
    url: String,
    loading: String,
  },
  computed: {
    host () {
      const parsedUrl = new URL(this.url)
      if (parsedUrl.hostname.startsWith('www.')) {
        return parsedUrl.hostname.substring(4)
      }
      return parsedUrl.hostname
    },
    src () {
      const parsedUrl = new URL(this.url)
      if (this.host === 'streamable.com') {
        return `https://streamable.com/o${parsedUrl.pathname}`
      } else if (this.host === 'streamja.com') {
        return `https://streamja.com/embed${parsedUrl.pathname}`
      } else if (this.host === 'clippituser.tv') {
        return `https://clippituser.tv/c/embed_iframe${parsedUrl.pathname.substring(2)}`
      }
      return null
    },
    hasEmbedFrame () {
      return this.src !== null
    },
  },
  methods: {
    loadFrame () {
      if (this.hasEmbedFrame && !this.$refs.frame.src) {
        this.$refs.frame.src = this.src
      }
      if (!this.hasEmbedFrame) {
        this.$refs.video.style.display = 'block'
      }
    },
  },
  mounted () {
    if (this.loading === 'eager') {
      this.loadFrame()
    } else if (this.loading === 'lazy') {
      this.$refs.frame.addEventListener('enteredViewport', this.loadFrame)
    }
  },
  beforeDestroy () {
    this.$refs.frame.removeEventListener('enteredViewport', this.loadFrame)
  },
}
</script>

<style scoped lang="less">
  .highlight-video {
    background-color: #ccc;

    &__frame {
      position: absolute;
      width: 100%;
      height: 100%;
    }
  }
</style>
