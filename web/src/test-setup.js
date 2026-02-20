// styled-components 6.x の CJS バンドルがグローバル React を参照するため設定
const React = require("react");
globalThis.React = React;
