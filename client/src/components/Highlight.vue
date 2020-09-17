<template>
  <div class="highlight" :class="{ 'highlight--smaller': isSmaller }">
    <highlight-video class="highlight__video" :loading="loading" :url="url"></highlight-video>
    <div class="highlight__meta">
      <div v-if="title" class="highlight__meta__title">
        <router-link :to="{ name: 'highlight', params: { id: id } }">
          {{ title }}
          <template v-if="numMirrors > 0">
            <span class="middot">&middot;</span>
            <span class="highlight__meta__title__mirrors">{{ numMirrors }} mirror<template v-if="numMirrors > 1">s</template></span>
          </template>
        </router-link>
      </div>
      <div class="highlight__meta__subtitle">
        <div class="highlight__meta__subtitle__submitted">
          Submitted {{ redditTimeAgo }} ago by <a :href="redditAuthorUrl">{{ redditAuthor }}</a>
        </div>
        <div class="highlight__meta__subtitle__comments">
          <span class="middot">&middot;</span><a :href="redditCommentsUrl">Comments</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { formatDistanceToNow, parseISO } from 'date-fns'
import HighlightVideo from './HighlightVideo.vue'

export default {
  name: 'Highlight',
  components: {
      HighlightVideo,
  },
  props: {
    id: String,
    loading: { type: String, default: 'eager' },
    title: String,
    url: String,
    redditPermalink: String,
    redditAuthor: String,
    redditCreatedAt: String,
    numMirrors: Number,
    isSmaller: { type: Boolean, default: false },
  },
  computed: {
    redditTimeAgo () {
      return formatDistanceToNow(parseISO(this.redditCreatedAt))
    },
    redditAuthorUrl () {
      return `https://reddit.com/u/${this.redditAuthor}`
    },
    redditCommentsUrl () {
      return `https://reddit.com${this.redditPermalink}`
    },
  },
}
</script>

<style scoped lang="less">
.middot {
  margin: 0 4px 0 4px;
}

.highlight {
    margin: 0 16px 16px 0;
    width: 600px;
    background-color: white;
    box-shadow: 0 1px 2px rgba(0,0,0,0.07),
                0 2px 4px rgba(0,0,0,0.07),
                0 4px 8px rgba(0,0,0,0.07),
                0 8px 16px rgba(0,0,0,0.07),
                0 16px 32px rgba(0,0,0,0.07),
                0 32px 64px rgba(0,0,0,0.07);
    border: 1px solid #ececf6;
    border-bottom-left-radius: 5px;
    border-bottom-right-radius: 5px;

    &--smaller {
      width: 350px;
    }

    &__meta {
      padding: 8px;
      border-bottom-left-radius: 5px;
      border-bottom-right-radius: 5px;

      &__title {
        font-size: 16px;
        margin-bottom: 8px;

        &__mirrors {
          text-decoration: underline;
        }

        a, a:visited {
          color: inherit;
          text-decoration: none;
        }
      }

      &__subtitle {
        display: flex;
        flex-wrap: wrap;
        font-size: 12px;
        color: #737373;

        a, a:visited {
          color: #737373;
        }
      }
    }

    &__video {
      width: 100%;
    }
}

@media only screen and (max-width: 600px) {
  .highlight {
    width: calc(100% - 16px);
    min-width: 350px;
    margin-right: 0;
  }
}
</style>
