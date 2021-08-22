#include <jni.h>

// See
// https://stackoverflow.com/questions/53252803/how-do-i-call-a-java-native-interface-c-function-from-my-go-code
const char *CGoGetStringUTFChars(JNIEnv *env, jstring str) {
  return (*env)->GetStringUTFChars(env, str, 0);
}