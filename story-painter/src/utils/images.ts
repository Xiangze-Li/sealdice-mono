import { trim } from "lodash-es";

export function extractImageUrl(text: string | undefined): string | undefined {
  if (!text || !trim(text)) {
    return undefined
  }

  // ob11 - gocq
  // 因为与新napcat的格式冲突，所以加了一点匹配内容
  let m = /url=(https?\:\/\/gchat\.qpic\.cn[^\]]+)]/.exec(text) as RegExpExecArray
  if (m && !text.includes('file_unique')) {
    return m[1]
  }

  m = /\[(?:image|图):([^\]]+)?([^\]]+)\]/.exec(text) as RegExpExecArray
  if (m) {
    return m[1]
  }

  // ob11 - napcat 20250317
  if (!m) {
    m = /file_unique=([a-zA-Z0-9]{32})\]/.exec(text) as RegExpExecArray
    if (m) {
      // 注: 观察到与ob11冲突情况
      const url = `https://gchat.qpic.cn/gchatpic_new/0/0-0-${m[1]}/0?term=2,subType=1`;
      return url
    }
  }

  // ob11 - lagrange
  m = /file=(https?:\/\/[^\]]+)\]/.exec(text) as RegExpExecArray
  if (m) {
    return m[1]
  }

  // ob11 - llob(new2)
  m = /file=\{([A-Z0-9]+)-([A-Z0-9]+)-([A-Z0-9]+)-([A-Z0-9]+)-([A-Z0-9]+)}([^\]]+?)\]/.exec(text) as RegExpExecArray
  if (m) {
    const url = `https://gchat.qpic.cn/gchatpic_new/0/0-0-${m[1]}${m[2]}${m[3]}${m[4]}${m[5]}/0?term=2,subType=1`
    return url
  }

  // ob11 - llob(new)
  m = /file=([A-Za-z0-9]{32,64})(\.[a-zA-Z]+?)\]/.exec(text) as RegExpExecArray
  if (m) {
    const url = `https://gchat.qpic.cn/gchatpic_new/0/0-0-${m[1]}/0?term=2,subType=1`;
    return url
  }

  // ob11 - llob(old)
  m = /file=file:\/\/[^\]]+([A-Za-z0-9]{32})(\.[a-zA-Z]+?)\]/.exec(text) as RegExpExecArray
  if (m) {
    const url = `https://gchat.qpic.cn/gchatpic_new/0/0-0-${m[1].toUpperCase()}/0?term=2,subType=1`;
    return url
  }

  m = /\[mirai:image:\{([A-Z0-9]+)-([A-Z0-9]+)-([A-Z0-9]+)-([A-Z0-9]+)-([A-Z0-9]+)}([^\]]+?)\]/.exec(text) as RegExpExecArray
  if (m) {
    const url = `https://gchat.qpic.cn/gchatpic_new/0/0-0-${m[1]}${m[2]}${m[3]}${m[4]}${m[5]}/0?term=2`
    return url
  }
  return undefined
}

