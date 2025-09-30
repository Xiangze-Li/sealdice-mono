import { type FunctionalComponent } from 'vue';

export interface SensitiveTagProps {
  type: 'default' | 'info' | 'warning' | 'error';
  message?: string;
}

const SensitiveTag: FunctionalComponent<SensitiveTagProps> = props => {
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
  return (
    <n-tag size='small' type={props.type} bordered={false}>
      {msg}
    </n-tag>
  );
};

export default SensitiveTag;
