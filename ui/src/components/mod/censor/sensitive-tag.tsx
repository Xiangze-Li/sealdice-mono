import { defineComponent } from 'vue';
import { NTag } from 'naive-ui';

export interface SensitiveTagProps {
  type: 'default' | 'info' | 'warning' | 'error';
  message?: string;
}

export default defineComponent({
  name: 'SensitiveTag',
  props: {
    type: {
      type: String as PropType<SensitiveTagProps['type']>,
    },
    message: {
      type: String as PropType<SensitiveTagProps['message']>,
    },
  },
  setup(props: SensitiveTagProps) {
    let msg = props.message;
    if (!msg) {
      if (props.type === 'default') {
        msg = '提醒';
      } else if (props.type === 'info') {
        msg = '注意';
      } else if (props.type === 'warning') {
        msg = '警告';
      } else if (props.type === 'error') {
        msg = '危险';
      }
    }
    return () => {
      return (
        <NTag size="small" type={props.type} bordered={false}>
          {msg}
        </NTag>
      );
    };
  },
});
