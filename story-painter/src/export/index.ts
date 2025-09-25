import { LogItem, RenderConfig } from "~/types";
import { AlignmentType, Document, ImageRun, Packer, Paragraph, TextRun } from "docx";
import dayjs from "dayjs";
import { trim } from "es-toolkit/compat";
import { extractImageUrl } from "~/utils/images";
import axios from "axios";
import { imageSize } from "image-size";
import { msgImageFormat } from "~/utils";

const CREATOR = "SealDice Story Painter";

export async function toDocx(items: LogItem[], config: RenderConfig): Promise<Blob> {
  console.log('export to docs: item =', items, ', config =', config)
  let result: Paragraph[] = []
  for (let i of items) {
    let p = await renderItem(i, config)
    result.push(p)
  }

  const docx = new Document({
    creator: CREATOR,
    sections: [
      {
        properties: {},
        children: result,
      }
    ],
  });
  return Packer.toBlob(docx)
}

async function renderItem(item: LogItem, config: RenderConfig): Promise<Paragraph> {
  const timeFormat = config.yearHide ? 'HH:mm:ss' : 'YYYY/MM/DD HH:mm:ss'
  const imgUrl = extractImageUrl(item.message)
  const message = msgImageFormat(item.message, { imageHide: true })

  const elems = []

  if (!config.timeHide) {
    elems.push(
      new TextRun({ text: item.time ? `${dayjs.unix(item.time).format(timeFormat)} ` : '', color: '#AAAAAA' }),
    )
  }

  if (item.nickname) {
    elems.push(
      new TextRun({ text: `\<${item.nickname}\>:`, color: item.color }),
    )
  }

  if (!config.imageHide && imgUrl) {
    try {
      const imgUrlProxy = `https://cors-anywhere.com/${imgUrl}` // 为了解决跨域问题
      const response = await axios.get(imgUrlProxy, { responseType: 'arraybuffer' });
      const data = response.data;
      console.log('download image, url:', imgUrlProxy, ', data:', data)
      const { width: rawW, height: rawH } = imageSize(new Uint8Array(data))
      const { width, height } = calculateAspectRatioFit(rawW, rawH, 300, 300)
      elems.push(
        new ImageRun({
          type: 'jpg',
          data: data,
          transformation: {
            width: width,
            height: height,
          },
        }),
      )
    } catch (error) {
      console.error('Error downloading image:', error);
      elems.push(
        new TextRun({ text: '[未知图片]', color: item.color }),
      )
    }
  }

  if (message) {
    const lines = message?.trim()?.split('\n') ?? []
    for (let i = 0; i < lines.length; i++){
      let line = lines[i];
      elems.push(
        new TextRun({ text: line, color: item.color, break: i != 0 ? 1 : 0 })
      )
    }
  }

  return new Paragraph({
    children: elems,
    alignment: AlignmentType.LEFT,
  })
}

/**
 * 计算图片缩放
 */
function calculateAspectRatioFit(srcWidth: number, srcHeight: number, maxWidth: number, maxHeight: number): { width: number, height: number } {
  const ratio = Math.min(maxWidth / srcWidth, maxHeight / srcHeight);
  // 如果图片本身就比限制要小，则使用原始尺寸
  if (ratio >= 1) {
    return { width: srcWidth, height: srcHeight };
  }
  return { width: srcWidth * ratio, height: srcHeight * ratio };
}