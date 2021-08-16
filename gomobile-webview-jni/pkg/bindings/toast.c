
#include <jni.h>

void show_toast(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx) {
  JavaVM *vm = (JavaVM *)java_vm;
  JNIEnv *env = (JNIEnv *)jni_env;
  jobject context = (jobject)ctx;

  jclass system = (*env)->FindClass(env, "java/lang/System");
  jfieldID id =
      (*env)->GetStaticFieldID(env, system, "out", "Ljava/io/PrintStream;");
  jobject obj = (*env)->GetStaticObjectField(env, system, id);
  jclass cls = (*env)->GetObjectClass(env, obj);
  jmethodID println =
      (*env)->GetMethodID(env, cls, "println", "(Ljava/lang/String;)V");
  jobject hello = (*env)->NewStringUTF(env, "Hello, world!");
  (*env)->CallVoidMethod(env, obj, println, hello);
}