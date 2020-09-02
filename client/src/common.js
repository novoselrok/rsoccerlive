import { format, parse } from 'date-fns'

const DAY_FORMAT = 'yyyy-MM-dd'

export const formatDay = day => format(day, DAY_FORMAT)
export const parseDay = day => parse(day, DAY_FORMAT, new Date())

export const fetchApi = endpoint => fetch(`${process.env.VUE_APP_RSOCCERLIVE_API_URL}/api${endpoint}`).then(response => response.json())
