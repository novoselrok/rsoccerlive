import Vue from 'vue'
import Vuex from 'vuex'
import VCalendar from 'v-calendar';
import { formatDay, fetchApi } from '@/common.js'

Vue.use(Vuex)
Vue.use(VCalendar)

export default new Vuex.Store({
  state: {
    dayHighlights: {},
  },
  mutations: {
    SET_DAY_HIGHLIGHTS (state, { day, highlights }) {
      Vue.set(state.dayHighlights, day, highlights)
    },
  },
  actions: {
    fetchDayHighlights ({ state, commit }, date) {
      const dayFormatted = formatDay(date)
      if (state.dayHighlights[dayFormatted]) {
        return
      }
      fetchApi(`/highlights?day=${dayFormatted}`)
        .then(highlights => {
          commit('SET_DAY_HIGHLIGHTS', { day: dayFormatted, highlights })
        })
    },
  },
})
