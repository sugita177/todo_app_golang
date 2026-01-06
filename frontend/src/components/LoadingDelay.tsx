import { useState, useEffect, type ReactNode } from 'react';

interface Props {
  children: ReactNode;
  delay?: number; // ミリ秒単位
}

export const LoadingDelay = ({ children, delay = 300 }: Props) => {
  const [shouldShow, setShouldShow] = useState(false);

  useEffect(() => {
    // delay時間後に表示フラグを立てる
    const timer = setTimeout(() => {
      setShouldShow(true);
    }, delay);

    // クリーンアップ（コンポーネントが消えたらタイマー解除）
    return () => clearTimeout(timer);
  }, [delay]);

  // shouldShowがtrueになるまでは何も表示しない（nullを返す）
  return shouldShow ? <>{children}</> : null;
};