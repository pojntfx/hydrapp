#include "main.h"

jstring get_java_string(JNIEnv *env, char *msg) {
  return (*env)->NewStringUTF(env, msg);
}