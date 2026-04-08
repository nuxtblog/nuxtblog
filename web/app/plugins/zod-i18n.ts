import { z } from "zod";

const MESSAGES: Record<string, Record<string, string>> = {
  zh: {
    required: "此项为必填",
    invalid_type: "类型不正确",
    string_min: "至少需要 {n} 个字符",
    string_max: "最多 {n} 个字符",
    number_min: "不能小于 {n}",
    number_max: "不能大于 {n}",
    array_min: "至少需要 {n} 项",
    array_max: "最多 {n} 项",
    too_small: "值太小",
    too_big: "值太大",
    email: "请输入有效的邮箱地址",
    url: "请输入有效的 URL",
    pattern: "格式不正确",
    invalid_string: "格式不正确",
    invalid_option: "请选择有效选项",
  },
};

function msg(lang: string, key: string, vars?: Record<string, number>): string | undefined {
  const map = MESSAGES[lang];
  if (!map) return undefined;
  let s = map[key];
  if (!s) return undefined;
  if (vars) {
    for (const [k, v] of Object.entries(vars)) {
      s = s.replace(`{${k}}`, String(v));
    }
  }
  return s;
}

export default defineNuxtPlugin((nuxtApp) => {
  const i18n = nuxtApp.$i18n as { locale: { value: string } };

  function applyZodLocale(lang: string) {
    z.config({
      customError(iss) {
        switch (iss.code) {
          case "invalid_type":
            if (iss.received === "undefined" || iss.received === "null")
              return msg(lang, "required");
            return msg(lang, "invalid_type");

          case "too_small":
            if (iss.type === "string")
              return iss.minimum === 1
                ? msg(lang, "required")
                : msg(lang, "string_min", { n: Number(iss.minimum) });
            if (iss.type === "number")
              return msg(lang, "number_min", { n: Number(iss.minimum) });
            if (iss.type === "array")
              return msg(lang, "array_min", { n: Number(iss.minimum) });
            return msg(lang, "too_small");

          case "too_big":
            if (iss.type === "string")
              return msg(lang, "string_max", { n: Number(iss.maximum) });
            if (iss.type === "number")
              return msg(lang, "number_max", { n: Number(iss.maximum) });
            if (iss.type === "array")
              return msg(lang, "array_max", { n: Number(iss.maximum) });
            return msg(lang, "too_big");

          case "invalid_string":
            if (iss.validation === "email") return msg(lang, "email");
            if (iss.validation === "url") return msg(lang, "url");
            if (iss.validation === "regex") return msg(lang, "pattern");
            return msg(lang, "invalid_string");

          case "invalid_enum_value":
            return msg(lang, "invalid_option");
        }
      },
    });
  }

  applyZodLocale(i18n.locale.value);

  watch(
    () => i18n.locale.value,
    (lang) => applyZodLocale(lang),
  );
});
