#include "main.h"
#include "android/native_activity.h"

void ANativeActivity_onCreate(ANativeActivity *activity, void *savedState,
                              size_t savedStateSize) {
  // Cast the Java env and context
  JNIEnv *env = activity->env;
  jobject context = activity->clazz;

  // Get the toast functions
  jclass toast_cls = (*env)->FindClass(env, "android/widget/Toast");
  jmethodID toast_cls_make_text =
      (*env)->GetStaticMethodID(env, toast_cls, "makeText",
                                "(Landroid/content/Context;Ljava/lang/"
                                "CharSequence;I)Landroid/widget/Toast;");

  // Get the message to display
  jobject msg_obj = (*env)->NewStringUTF(env, "Hello from C!");

  // Create the toast
  jobject toast_obj = (*env)->CallStaticObjectMethod(
      env, toast_cls, toast_cls_make_text, context, msg_obj, 0);

  // Show the toast
  jclass toast_obj_class = (*env)->GetObjectClass(env, toast_obj);
  jmethodID toast_obj_class_show =
      (*env)->GetMethodID(env, toast_obj_class, "show", "()V");
  (*env)->CallVoidMethod(env, toast_obj, toast_obj_class_show);

  // Start the main loop
  GoLoop(activity);
}