#include <jni.h>

void show_toast(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx,
                char *msg) {
  // Cast the Java VM, env and context
  JavaVM *vm = (JavaVM *)java_vm;
  JNIEnv *env = (JNIEnv *)jni_env;
  jobject context = (jobject)ctx;

  // Get the looper functions
  jclass looper_cls = (*env)->FindClass(env, "android/os/Looper");
  jmethodID looper_cls_prepare =
      (*env)->GetStaticMethodID(env, looper_cls, "prepare", "()V");
  jmethodID lopper_cls_loop =
      (*env)->GetStaticMethodID(env, looper_cls, "loop", "()V");

  // Get the toast functions
  jclass toast_cls = (*env)->FindClass(env, "android/widget/Toast");
  jmethodID toast_cls_make_text =
      (*env)->GetStaticMethodID(env, toast_cls, "makeText",
                                "(Landroid/content/Context;Ljava/lang/"
                                "CharSequence;I)Landroid/widget/Toast;");

  // Lock the UI thread
  (*env)->CallStaticVoidMethod(env, looper_cls, looper_cls_prepare);

  // Get the message to display
  jobject msg_obj = (*env)->NewStringUTF(env, msg);

  // Create the toast
  jobject toast_obj = (*env)->CallStaticObjectMethod(
      env, toast_cls, toast_cls_make_text, context, msg_obj, 0);

  // Show the toast
  jclass toast_obj_class = (*env)->GetObjectClass(env, toast_obj);
  jmethodID toast_obj_class_show =
      (*env)->GetMethodID(env, toast_obj_class, "show", "()V");
  (*env)->CallVoidMethod(env, toast_obj, toast_obj_class_show);

  // Unlock the UI thread
  //   (*env)->CallStaticVoidMethod(env, looper_cls, lopper_cls_loop);
}